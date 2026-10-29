$(function () {
  const error = function (message) { $("#fuelError").text(message).prop("hidden", false); };
  let socket;
  const phaseNames = ["PRE-MATCH", "STARTING", "AUTO", "PAUSE", "TELEOP", "POST-MATCH", "TIMEOUT", "POST-TIMEOUT"];
  const storage = {getItem: key => window.sessionStorage.getItem(key), setItem: (key, value) => window.sessionStorage.setItem(key, value)};
  try { $("#fuelOperator").val(storage.getItem("fuelOperator") || ""); }
  catch (_) { error("Browser storage is unavailable; pending actions cannot be saved."); }
  $("#fuelOperator").on("input", function () {
    try { storage.setItem("fuelOperator", this.value); }
    catch (_) { error("Cannot save the operator label in this browser."); }
  });
  const client = new ManualFuelClient({
    storage,
    uuid: () => typeof crypto.randomUUID === "function" ? crypto.randomUUID() : "10000000-1000-4000-8000-100000000000".replace(/[018]/g, c => (Number(c) ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> Number(c) / 4).toString(16)),
    operator: () => $("#fuelOperator").val(),
    send: (kind, command) => socket.send(kind, command),
    error,
    render: function (state, connected, pending) {
      $("#fuelConnection").text(connected ? "● Connected" : "○ Connecting / disconnected").toggleClass("text-success", connected);
      $("#fuelPending").text(pending + " pending");
      $("[data-delta]").prop("disabled", true);
      $("#fuelMode").prop("disabled", true);
      if (!state) return;
      $("#fuelMatch").text(state.MatchName);
      $("#fuelPhase").text(phaseNames[state.Phase]);
      $("#fuelRecovery").prop("hidden", !state.Recovered);
      $("#fuelMode").val(state.Mode).prop("disabled", !connected || !state.CanSetMode || pending > 0);
      $("#fuelSaved").text(state.CheckpointAt && !state.CheckpointAt.startsWith("0001") ? "Checkpoint: " + new Date(state.CheckpointAt).toLocaleTimeString() : "Awaiting checkpoint");
      if (state.StorageError) error("Storage error: " + state.StorageError);
      for (const alliance of ["red", "blue"]) {
        const values = state[alliance === "red" ? "Red" : "Blue"];
        const manual = state.Mode === "manual";
        const source = values.Source === "arduino" ? "Arduino" : "PLC";
        $("#" + alliance + "FuelSource").text(manual ? "Manual FUEL · other PLC functions remain enabled" : source + " base + operator corrections");
        $("#" + alliance + "FuelWarning").text(values.Suspended).prop("hidden", manual || !values.Suspended);
        $("[data-rebase='" + alliance + "']").prop("hidden", manual || source !== "PLC" || !values.Suspended).prop("disabled", !connected || !state.CanRebase || pending > 0);
        for (const phase of ["auto", "teleop"]) {
          const key = phase === "auto" ? "Auto" : "Teleop";
          const count = values[key];
          $("#" + alliance + "-" + phase + "-value").text(count);
          $("#" + alliance + "-" + phase + "-source").text(manual ? "Manual scored FUEL" : source + " " + values["Raw" + key] + " · correction " + (values[key + "Adjustment"] >= 0 ? "+" : "") + values[key + "Adjustment"]);
          $(".fuel-counter[data-alliance='" + alliance + "'][data-phase='" + phase + "'] [data-delta]").each(function () {
            this.disabled = !connected || !state["Can" + key] || pending >= 32 || count + Number(this.dataset.delta) < 0;
          });
        }
      }
    }
  });
  $("[data-delta]").on("click", function () {
    $("#fuelError").prop("hidden", true);
    const counter = $(this).closest(".fuel-counter");
    client.action("adjustFuel", {Alliance: counter.data("alliance"), Phase: counter.data("phase"), Delta: Number(this.dataset.delta)});
  });
  $("#fuelMode").on("change", function () { client.action("setFuelMode", {Mode: this.value}); });
  $("[data-rebase]").on("click", function () {
    if (confirm("Use the current PLC counter as the new baseline? No FUEL will be added. Correct any missed FUEL manually.")) client.action("rebasePlcFuel", {Alliance: this.dataset.rebase});
  });
  $("#fuelHistoryRefresh").on("click", function () {
    if (!client.state) return;
    $.getJSON("/panels/manual_fuel/history", {attemptId: client.state.HistoryId || client.state.AttemptId}).done(function (rows) {
      const body = $("#fuelHistory").empty();
      for (const row of rows) {
        const tr = $("<tr>");
        for (const value of [new Date(row.Time).toLocaleTimeString(), row.Command.Operator || row.Actor, [row.Kind, row.Command.Alliance, row.Command.Phase, row.Command.Delta || ""].join(" "), row.Before + " → " + row.After]) tr.append($("<td>").text(value));
        body.append(tr);
      }
    }).fail(() => error("Cannot load the audit log."));
  });
  socket = new CheesyWebsocket("/panels/manual_fuel/websocket", {
    manualFuelState: event => client.snapshot(event.data),
    fuelResult: event => client.result(event.data),
    error: event => error(event.data)
  }, {open: () => client.open(), close: () => client.close()});
});
