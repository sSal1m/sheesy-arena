$(function () {
  const error = message => $("#arduinoError").text(message).prop("hidden",false);
  function refresh() {
    $.getJSON("/panels/arduino/status").done(function (status) {
      $("#arduinoListener").text(status.Enabled ? "Authenticated Ethernet listener: " + status.Listen : "Arduino integration disabled. Configure the device file and restart before a match.");
      const root = $("#arduinoDevices").empty();
      for (const device of status.Devices) {
        const card = $("<section class='col-lg-6'><div class='card h-100'><div class='card-body'></div></div></section>");
        const body = card.find(".card-body");
        card.find(".card").addClass(device.Alliance === "red" ? "border-danger" : "border-primary");
        $("<h2 class='h4'>").text(device.Alliance.toUpperCase() + " · " + device.Id).appendTo(body);
        $("<p class='fw-semibold'>").toggleClass("text-success",device.Connected).text(device.Connected ? "Connected" : "Disconnected").appendTo(body);
        const rows = [["IP", device.IP], ["Current / configured automatic source",device.ActiveSource + " / " + device.Source], ["Last communication", device.LastSeen.startsWith("0001") ? "Never" : new Date(device.LastSeen).toLocaleString()], ["Firmware / boot",device.Firmware + " / " + device.Boot],["Detected / accepted",device.Detected + " / " + device.Accepted],["Rejected / duplicates / device drops",device.Rejected + " / " + device.Duplicates + " / " + device.Dropped],["Last error",device.LastError || "—"]];
        const state = device.State;
        if (state) rows.push(["Match / phase", state.MatchId + " / " + ["PRE-MATCH","STARTING","AUTO","PAUSE","TELEOP","POST-MATCH","TIMEOUT","POST-TIMEOUT"][state.Phase]], ["HUB / scoring", (state.Flags & 8 ? "active" : "inactive") + " / " + (state.Flags & 16 ? "enabled" : "disabled")]);
        for (const [label,value] of rows) $("<p class='mb-2'>").append($("<span class='text-body-secondary'>").text(label + ": "),$("<span>").text(value)).appendTo(body);
        root.append(card);
      }
    }).fail(() => error("Cannot read Arduino status; check the admin session and server connection."));
  }
  function simulate(alliance,reset) {
    $.ajax({url:"/panels/arduino/simulate",method:"POST",contentType:"application/json",data:JSON.stringify({Alliance:alliance,Phase:$("#simulationPhase").val(),Count:reset ? 0 : 1,Reset:reset})}).done(result => {
      $("#simulationResult").text("SIMULATION ONLY · RED: AUTO " + result.Red.Auto + ", TELEOP " + result.Red.Teleop + ", total " + result.RedSummary.Score + " · BLUE: AUTO " + result.Blue.Auto + ", TELEOP " + result.Blue.Teleop + ", total " + result.BlueSummary.Score);
    }).fail(xhr => error("Simulation failed: " + (xhr.responseText || "connection error")));
  }
  $("[data-simulate]").on("click",function(){simulate(this.dataset.simulate,false);});
  $("#simulationReset").on("click",()=>simulate("red",true));
  refresh(); setInterval(refresh,1000);
});
