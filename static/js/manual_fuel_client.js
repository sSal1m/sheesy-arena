// Transport-independent state/outbox logic, shared by the UI and Node tests.
(function (root) {
  class ManualFuelClient {
    constructor(options) {
      this.options = options;
      this.state = null;
      this.connected = false;
      this.awaitingSnapshot = true;
      try { this.pending = JSON.parse(options.storage.getItem("fuelOutbox") || "[]"); }
      catch (_) { this.pending = []; }
      if (!Array.isArray(this.pending)) this.pending = [];
      this.pending = this.pending.filter(item => item && typeof item.kind === "string" && item.command && typeof item.command.RequestId === "string" && typeof item.command.AttemptId === "string").slice(0, 32);
    }
    persist() { this.options.storage.setItem("fuelOutbox", JSON.stringify(this.pending)); }
    persistSafely() {
      try { this.persist(); }
      catch (_) { this.options.error("Cannot save pending actions in this browser; check the server state before retrying."); }
    }
    open() { this.connected = true; this.awaitingSnapshot = true; this.render(); }
    close() { this.connected = false; this.render(); }
    render() { this.options.render(this.state, this.connected && !this.awaitingSnapshot, this.pending.length); }
    snapshot(state) {
      const first = this.awaitingSnapshot;
      if (!first && this.state && state.Revision < this.state.Revision) return;
      this.awaitingSnapshot = false;
      this.state = state;
      if (first) {
        const retained = [];
        for (const item of this.pending) {
          if (item.command.AttemptId === state.AttemptId && item.command.MatchId === state.MatchId) {
            retained.push(item);
            try { this.options.send(item.kind, item.command); }
            catch (_) { this.options.error("Connection changed; the pending action will be retried after reconnect."); }
          } else this.options.error("A pending action belonged to an earlier match context. It was not replayed; check the audit log.");
        }
        this.pending = retained;
        this.persistSafely();
      }
      this.render();
    }
    result(reply) {
      this.pending = this.pending.filter(item => item.command.RequestId !== reply.RequestId);
      this.persistSafely();
      if (reply.Error) this.options.error(reply.Error);
      this.snapshot(reply.State);
    }
    action(kind, fields) {
      if (!this.connected || this.awaitingSnapshot || !this.state) return;
      if (this.pending.length >= 32) { this.options.error("Too many pending actions; wait for server confirmation."); return; }
      const command = Object.assign({RequestId: this.options.uuid(), AttemptId: this.state.AttemptId, MatchId: this.state.MatchId, Operator: this.options.operator()}, fields);
      const item = {kind, command};
      this.pending.push(item);
      try { this.persist(); }
      catch (_) { this.pending.pop(); this.options.error("Cannot save the pending action in this browser. No change was sent."); return; }
      try { this.options.send(kind, command); }
      catch (_) { this.options.error("Connection changed; the pending action will be retried after reconnect."); }
      this.render();
    }
  }
  root.ManualFuelClient = ManualFuelClient;
  if (typeof module !== "undefined") module.exports = ManualFuelClient;
})(typeof window === "undefined" ? globalThis : window);
