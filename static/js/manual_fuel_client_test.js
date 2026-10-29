const test = require("node:test");
const assert = require("node:assert/strict");
const Client = require("./manual_fuel_client.js");

function setup(existing = []) {
  const values = new Map([["fuelOutbox", JSON.stringify(existing)]]);
  const sent = [], errors = [], renders = [];
  let id = 0;
  const options = {storage: {getItem: key => values.get(key), setItem: (key, value) => values.set(key, value)}, uuid: () => "id-" + ++id, operator: () => "A", send: (kind, command) => sent.push({kind, command}), error: value => errors.push(value), render: (...args) => renders.push(args)};
  return {client: new Client(options), options, sent, errors, renders};
}
const state = (revision = 1, attempt = "a") => ({Revision: revision, AttemptId: attempt, MatchId: 1, Red: {Auto: revision}});

test("values come only from authoritative snapshots, not optimistic increments", () => {
  const {client, sent} = setup();
  client.open(); client.snapshot(state());
  client.action("adjustFuel", {Alliance: "red", Phase: "auto", Delta: 1});
  assert.equal(client.state.Red.Auto, 1);
  assert.equal(sent.length, 1);
  client.result({RequestId: sent[0].command.RequestId, State: state(2)});
  assert.equal(client.state.Red.Auto, 2);
  assert.equal(client.pending.length, 0);
});
test("disconnect/reconnect retries exactly the same request ID", () => {
  const {client, sent} = setup();
  client.open(); client.snapshot(state()); client.action("adjustFuel", {Delta: 1});
  client.close(); client.action("adjustFuel", {Delta: 1});
  assert.equal(sent.length, 1);
  client.open(); client.snapshot(state(2));
  assert.equal(sent.length, 2);
  assert.equal(sent[0].command.RequestId, sent[1].command.RequestId);
});
test("page refresh retains uncertain requests, not local score totals", () => {
  const context = setup();
  context.client.open(); context.client.snapshot(state()); context.client.action("adjustFuel", {Delta: 1});
  const reloaded = new Client(context.options);
  assert.equal(reloaded.state, null);
  reloaded.open(); reloaded.snapshot(state(3));
  assert.equal(context.sent[0].command.RequestId, context.sent[1].command.RequestId);
});
test("stale responses do not rewind score and duplicate ACK clears pending", () => {
  const {client, sent} = setup();
  client.open(); client.snapshot(state(4)); client.action("adjustFuel", {Delta: 1});
  client.snapshot(state(6));
  client.result({RequestId: sent[0].command.RequestId, Duplicate: true, State: state(5)});
  assert.equal(client.state.Revision, 6);
  assert.equal(client.pending.length, 0);
});
test("new match context drops old outbox without replaying", () => {
  const {client, sent, errors} = setup();
  client.open(); client.snapshot(state()); client.action("adjustFuel", {Delta: 1});
  client.close(); client.open(); client.snapshot(state(10, "b"));
  assert.equal(sent.length, 1);
  assert.equal(client.pending.length, 0);
  assert.equal(errors.length, 1);
});
test("server errors are visible and never locally applied", () => {
  const {client, sent, errors} = setup();
  client.open(); client.snapshot(state()); client.action("adjustFuel", {Delta: -5});
  client.result({RequestId: sent[0].command.RequestId, Error: "negative", State: state()});
  assert.deepEqual(errors, ["negative"]);
  assert.equal(client.state.Red.Auto, 1);
});
test("new server snapshot is accepted after reconnect even at lower revision", () => {
  const {client} = setup();
  client.open(); client.snapshot(state(100)); client.close(); client.open(); client.snapshot(state(20, "recovered"));
  assert.equal(client.state.AttemptId, "recovered");
});

test("a transport send error preserves the durable outbox for retry", () => {
  const {client, options, errors} = setup();
  client.open(); client.snapshot(state());
  options.send = () => { throw new Error("closed"); };
  client.action("adjustFuel", {Delta: 1});
  assert.equal(client.pending.length, 1);
  assert.equal(errors.length, 1);
});

test("malformed stored outbox does not prevent an initial snapshot", () => {
  const {client} = setup([null, {}, {kind: "bad"}]);
  client.open(); client.snapshot(state());
  assert.equal(client.pending.length, 0);
});

test("unavailable browser storage still displays server state and prevents untracked sends", () => {
  const {client, options, errors, sent, renders} = setup();
  options.storage.setItem = () => { throw new Error("quota"); };
  client.open(); client.snapshot(state());
  assert.equal(client.state.Red.Auto, 1);
  assert.ok(renders.length > 1);
  client.action("adjustFuel", {Delta: 1});
  assert.equal(sent.length, 0);
  assert.equal(client.pending.length, 0);
  client.result({RequestId: "old", State: state(2)});
  assert.equal(client.state.Red.Auto, 2);
  assert.ok(errors.length >= 2);
});
