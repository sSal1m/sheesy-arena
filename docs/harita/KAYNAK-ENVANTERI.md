# Kaynak ve bağımlılık envanteri

Satır sayıları yorumları ve boş satırları içerir. Test fonksiyonu sayısı `func Test...` tanımlarıdır; alt testler sayılmaz. `*_test.go` dosyaları üretim bağımlılıklarına dahil edilmez; `field/test_helpers.go` ise adına rağmen normal derlemeye dahildir.

| Paket | Go dosyası | Go satırı | Test dosyası | Test fonksiyonu | İç bağımlılıklar |
| --- | --- | --- | --- | --- | --- |
| . | 1 | 33 | 0 | 0 | field, network, web |
| field | 11 | 3994 | 8 | 43 | game, led, model, network, partner, playoff, plc, websocket |
| game | 9 | 810 | 7 | 26 | — |
| led | 5 | 766 | 2 | 9 | — |
| model | 17 | 1596 | 15 | 50 | game |
| network | 3 | 615 | 3 | 5 | model |
| partner | 4 | 1151 | 4 | 16 | game, model |
| playoff | 9 | 1387 | 5 | 32 | game, model |
| plc | 5 | 700 | 2 | 16 | websocket |
| tournament | 6 | 819 | 5 | 16 | game, model |
| web | 37 | 7375 | 35 | 119 | field, game, led, model, partner, playoff, tournament, websocket |
| websocket | 2 | 253 | 2 | 3 | — |

## Dosya sorumlulukları ve fonksiyon indeksi

### [field/arena.go](../../field/arena.go)

Functions for controlling the arena and match play.

[NewArena](../../field/arena.go#L127), [(arena *Arena).LoadSettings](../../field/arena.go#L173), [(arena *Arena).CreatePlayoffTournament](../../field/arena.go#L308), [(arena *Arena).CreatePlayoffMatches](../../field/arena.go#L317), [(arena *Arena).UpdatePlayoffTournament](../../field/arena.go#L322), [(arena *Arena).LoadMatch](../../field/arena.go#L334), [(arena *Arena).LoadTestMatch](../../field/arena.go#L419), [(arena *Arena).LoadNextMatch](../../field/arena.go#L424), [(arena *Arena).SubstituteTeams](../../field/arena.go#L457), [(arena *Arena).StartMatch](../../field/arena.go#L510), [(arena *Arena).AbortMatch](../../field/arena.go#L551), [(arena *Arena).ResetMatch](../../field/arena.go#L573), [(arena *Arena).StartTimeout](../../field/arena.go#L593), [(arena *Arena).StartAdHocTimeout](../../field/arena.go#L598), [(arena *Arena).SetTimeoutDisplay](../../field/arena.go#L603), [(arena *Arena).startTimeout](../../field/arena.go#L609), [(arena *Arena).currentMatchDisplayName](../../field/arena.go#L634), [(arena *Arena).SetAudienceDisplayMode](../../field/arena.go#L646), [(arena *Arena).SetAllianceStationDisplayMode](../../field/arena.go#L664), [(arena *Arena).MatchTimeSec](../../field/arena.go#L676), [(arena *Arena).Update](../../field/arena.go#L686), [(arena *Arena).checkEndgameStart](../../field/arena.go#L837), [(arena *Arena).logTeamSnapshots](../../field/arena.go#L856), [(arena *Arena).Run](../../field/arena.go#L877), [(arena *Arena).RedScoreSummary](../../field/arena.go#L903), [(arena *Arena).BlueScoreSummary](../../field/arena.go#L908), [(arena *Arena).validateTeams](../../field/arena.go#L913), [(allianceStation *AllianceStation).startTeamMatchLog](../../field/arena.go#L930), [(allianceStation *AllianceStation).closeTeamMatchLog](../../field/arena.go#L949), [(arena *Arena).closeTeamMatchLogs](../../field/arena.go#L957), [(arena *Arena).assignTeam](../../field/arena.go#L964), [(arena *Arena).getNextMatch](../../field/arena.go#L1005), [(arena *Arena).preLoadNextMatch](../../field/arena.go#L1025), [(arena *Arena).setSCCEthernetEnabled](../../field/arena.go#L1064), [(arena *Arena).setupNetwork](../../field/arena.go#L1083), [(arena *Arena).checkCanStartMatch](../../field/arena.go#L1110), [(arena *Arena).getStartMatchConditions](../../field/arena.go#L1119), [(arena *Arena).checkAllianceStationsReady](../../field/arena.go#L1161), [(arena *Arena).getAllianceStationStartConditions](../../field/arena.go#L1169), [(arena *Arena).sendDsPacket](../../field/arena.go#L1214), [(arena *Arena).getAssignedAllianceStation](../../field/arena.go#L1234), [(arena *Arena).handleAutoWinner](../../field/arena.go#L1245), [(arena *Arena).handlePlcInputOutput](../../field/arena.go#L1268), [(arena *Arena).getHubLightStates](../../field/arena.go#L1358), [(arena *Arena).ToggleBypass](../../field/arena.go#L1398), [(arena *Arena).handleTeamStop](../../field/arena.go#L1407), [(arena *Arena).SignalVolunteers](../../field/arena.go#L1425), [(arena *Arena).SignalReset](../../field/arena.go#L1439), [(arena *Arena).handleSounds](../../field/arena.go#L1456), [(arena *Arena).PlaySound](../../field/arena.go#L1476), [(arena *Arena).positionPostMatchScoreReady](../../field/arena.go#L1482), [(arena *Arena).checkForUpdatedNexusLineup](../../field/arena.go#L1487), [(arena *Arena).runPeriodicTasks](../../field/arena.go#L1517), [(arena *Arena).AutomateAudienceDisplay](../../field/arena.go#L1524)

### [field/arena_leds.go](../../field/arena_leds.go)

Game-specific control of the 2026 Hub DMX lighting.

[(arena *Arena).updateHubLeds](../../field/arena_leds.go#L21), [(arena *Arena).updateTeleopHubLeds](../../field/arena_leds.go#L47)

### [field/arena_notifiers.go](../../field/arena_notifiers.go)

Contains configuration of the publish-subscribe notifiers that allow the arena to push updates to websocket clients.

[(arena *Arena).configureNotifiers](../../field/arena_notifiers.go#L48), [(arena *Arena).generateAllianceSelectionMessage](../../field/arena_notifiers.go#L72), [(arena *Arena).generateAllianceStationDisplayModeMessage](../../field/arena_notifiers.go#L86), [(arena *Arena).generateArenaStatusMessage](../../field/arena_notifiers.go#L90), [(arena *Arena).generateAudienceDisplayModeMessage](../../field/arena_notifiers.go#L123), [(arena *Arena).generateDisplayConfigurationMessage](../../field/arena_notifiers.go#L127), [(arena *Arena).generateEventStatusMessage](../../field/arena_notifiers.go#L138), [(arena *Arena).generateLowerThirdMessage](../../field/arena_notifiers.go#L142), [(arena *Arena).GenerateMatchLoadMessage](../../field/arena_notifiers.go#L149), [(arena *Arena).getAllianceYellowCard](../../field/arena_notifiers.go#L243), [(arena *Arena).generateMatchTimeMessage](../../field/arena_notifiers.go#L252), [(arena *Arena).generateMatchTimingMessage](../../field/arena_notifiers.go#L256), [(arena *Arena).generateRealtimeScoreMessage](../../field/arena_notifiers.go#L260), [(arena *Arena).GenerateScorePostedMessage](../../field/arena_notifiers.go#L281), [(arena *Arena).generateScoringStatusMessage](../../field/arena_notifiers.go#L388), [getAudienceAllianceScoreFields](../../field/arena_notifiers.go#L415), [getRulesViolated](../../field/arena_notifiers.go#L428)

### [field/display.go](../../field/display.go)

Model representing and methods for controlling a remote web display.

[DisplayFromUrl](../../field/display.go#L101), [(display *Display).ToUrl](../../field/display.go#L144), [(display *Display).generateDisplayConfigurationMessage](../../field/display.go#L169), [(arena *Arena).NextDisplayId](../../field/display.go#L174), [(arena *Arena).RegisterDisplay](../../field/display.go#L190), [(arena *Arena).UpdateDisplay](../../field/display.go#L224), [(arena *Arena).MarkDisplayDisconnected](../../field/display.go#L243), [(arena *Arena).purgeDisconnectedDisplays](../../field/display.go#L262)

### [field/driver_station_connection.go](../../field/driver_station_connection.go)

Model and methods for interacting with a team's Driver Station.

[driverStationTeamIdFromRemoteAddr](../../field/driver_station_connection.go#L76), [newDriverStationConnection](../../field/driver_station_connection.go#L92), [(arena *Arena).initializeUdpListener](../../field/driver_station_connection.go#L119), [(arena *Arena).listenForDsUdpPackets](../../field/driver_station_connection.go#L139), [(dsConn *DriverStationConnection).update](../../field/driver_station_connection.go#L210), [(dsConn *DriverStationConnection).close](../../field/driver_station_connection.go#L228), [(dsConn *DriverStationConnection).encodeControlPacket](../../field/driver_station_connection.go#L237), [(dsConn *DriverStationConnection).sendControlPacket](../../field/driver_station_connection.go#L338), [listenAddress](../../field/driver_station_connection.go#L356), [(arena *Arena).listenForDriverStations](../../field/driver_station_connection.go#L364), [(arena *Arena).serveDriverStations](../../field/driver_station_connection.go#L384), [readTaggedTcpPacket](../../field/driver_station_connection.go#L567), [(dsConn *DriverStationConnection).handleTcpConnection](../../field/driver_station_connection.go#L594), [(dsConn *DriverStationConnection).copyDsReportedStatus](../../field/driver_station_connection.go#L621), [(dsConn *DriverStationConnection).parseDsLogPacket](../../field/driver_station_connection.go#L630), [handleInvalidTcpConnection](../../field/driver_station_connection.go#L645), [closeTcpConn](../../field/driver_station_connection.go#L682), [(dsConn *DriverStationConnection).checkGameData](../../field/driver_station_connection.go#L688), [(dsConn *DriverStationConnection).sendGameDataPacketTcp](../../field/driver_station_connection.go#L707)

### [field/event_status.go](../../field/event_status.go)

Model and functions for reporting on event status.

[(arena *Arena).updateCycleTime](../../field/event_status.go#L26), [(arena *Arena).updateEarlyLateMessage](../../field/event_status.go#L61), [(arena *Arena).getEarlyLateMessage](../../field/event_status.go#L70)

### [field/realtime_score.go](../../field/realtime_score.go)

Model representing the current state of the score during a match.

[NewRealtimeScore](../../field/realtime_score.go#L19)

### [field/scoring_panel_registry.go](../../field/scoring_panel_registry.go)

Model representing and methods for tracking the state of a realtime scoring panel.

[(registry *ScoringPanelRegistry).initialize](../../field/scoring_panel_registry.go#L18), [(registry *ScoringPanelRegistry).resetScoreCommitted](../../field/scoring_panel_registry.go#L23), [(registry *ScoringPanelRegistry).GetNumPanels](../../field/scoring_panel_registry.go#L35), [(registry *ScoringPanelRegistry).GetNumScoreCommitted](../../field/scoring_panel_registry.go#L43), [(registry *ScoringPanelRegistry).RegisterPanel](../../field/scoring_panel_registry.go#L57), [(registry *ScoringPanelRegistry).SetScoreCommitted](../../field/scoring_panel_registry.go#L68), [(registry *ScoringPanelRegistry).UnregisterPanel](../../field/scoring_panel_registry.go#L76)

### [field/team_match_log.go](../../field/team_match_log.go)

Utilities for logging team station snapshots during a match.

[NewTeamMatchLog](../../field/team_match_log.go#L30), [(log *TeamMatchLog).LogStationSnapshot](../../field/team_match_log.go#L61), [(log *TeamMatchLog).Close](../../field/team_match_log.go#L141)

### [field/team_sign.go](../../field/team_sign.go)

Models and logic for controlling a Cypress team number / timer sign.

[NewTeamSigns](../../field/team_sign.go#L76), [(signs *TeamSigns).Update](../../field/team_sign.go#L84), [(signs *TeamSigns).SetNextMatchTeams](../../field/team_sign.go#L128), [(sign *TeamSign).SetId](../../field/team_sign.go#L138), [(sign *TeamSign).update](../../field/team_sign.go#L175), [generateInMatchTeamRearText](../../field/team_sign.go#L195), [generateInMatchTimerRearText](../../field/team_sign.go#L222), [generateTeamSignAllianceScores](../../field/team_sign.go#L228), [formatTeamSignRearText](../../field/team_sign.go#L248), [generateTeamSignPeriodText](../../field/team_sign.go#L253), [generateTeamSignPeriodPrefix](../../field/team_sign.go#L265), [generateTimerTexts](../../field/team_sign.go#L289), [(sign *TeamSign).generateTeamNumberTexts](../../field/team_sign.go#L327), [(sign *TeamSign).sendPacket](../../field/team_sign.go#L415), [(sign *TeamSign).writePacketData](../../field/team_sign.go#L460), [blinkColor](../../field/team_sign.go#L468)

### [field/test_helpers.go](../../field/test_helpers.go)

Helper methods for use in tests in this package and others.

[SetupTestArena](../../field/test_helpers.go#L17), [setupTestArena](../../field/test_helpers.go#L32)

### [game/foul.go](../../game/foul.go)

Model of a foul.

[(foul *Foul).Rule](../../game/foul.go#L16), [(foul *Foul).PointValue](../../game/foul.go#L21)

### [game/hub.go](../../game/hub.go)

Scoring logic for the 2026 Hub element.

[(hub *Hub).UpdateState](../../game/hub.go#L34), [(hub *Hub).GetTeleopActiveFuelCount](../../game/hub.go#L57), [(hub *Hub).GetShiftCount](../../game/hub.go#L67), [(hub *Hub).GetCurrentShiftTiming](../../game/hub.go#L76), [(hub *Hub).GetActiveShiftTiming](../../game/hub.go#L103), [(hub *Hub).isShiftActive](../../game/hub.go#L115), [(hub *Hub).getCurrentShift](../../game/hub.go#L130), [(hub *Hub).getScoringGracePeriod](../../game/hub.go#L159)

### [game/match_sounds.go](../../game/match_sounds.go)

Game-specific audience sound timings.

[UniqueMatchSounds](../../game/match_sounds.go#L19), [UpdateMatchSounds](../../game/match_sounds.go#L37)

### [game/match_timing.go](../../game/match_timing.go)

Game-specific period timing.

[GetTeleopDurationSec](../../game/match_timing.go#L24), [GetDurationToAutoEnd](../../game/match_timing.go#L28), [GetDurationToTeleopStart](../../game/match_timing.go#L32), [GetDurationToTeleopEnd](../../game/match_timing.go#L38)

### [game/ranking_fields.go](../../game/ranking_fields.go)

Game-specific fields by which teams are ranked and the logic for sorting rankings.

[GetWinRankingPoints](../../game/ranking_fields.go#L34), [(fields *RankingFields).AddScoreSummary](../../game/ranking_fields.go#L41), [(rankings Rankings).Len](../../game/ranking_fields.go#L72), [(rankings Rankings).Less](../../game/ranking_fields.go#L77), [(rankings Rankings).Swap](../../game/ranking_fields.go#L98)

### [game/rule.go](../../game/rule.go)

Model of a game-specific rule.

[GetRuleById](../../game/rule.go#L56), [GetAllRules](../../game/rule.go#L61)

### [game/score.go](../../game/score.go)

Model representing the instantaneous score of a match.

[(score *Score).Summarize](../../game/score.go#L32), [(score *Score).Equals](../../game/score.go#L123)

### [game/score_summary.go](../../game/score_summary.go)

Model representing the calculated totals of a match score.

[(t MatchStatus).Get](../../game/score_summary.go#L38), [DetermineMatchStatus](../../game/score_summary.go#L44), [comparePoints](../../game/score_summary.go#L83)

### [game/test_helpers.go](../../game/test_helpers.go)

Helper methods for use in tests in this package and others.

[TestScore1](../../game/test_helpers.go#L8), [TestScore2](../../game/test_helpers.go#L30), [TestRanking1](../../game/test_helpers.go#L43), [TestRanking2](../../game/test_helpers.go#L47)

### [led/color.go](../../led/color.go)

Contains pixel RGB mappings for common LED colors.

[(color Color).Scale](../../led/color.go#L24)

### [led/controller.go](../../led/controller.go)

Represents an E1.31 sACN (DMX over Ethernet) LED controller for the 2026 hub lights.

[NewController](../../led/controller.go#L40), [(controller *Controller).SetUniverseMode](../../led/controller.go#L50), [(controller *Controller).SetAddress](../../led/controller.go#L59), [(controller *Controller).SetMode](../../led/controller.go#L77), [(controller *Controller).GetModes](../../led/controller.go#L89), [(controller *Controller).GetPixels](../../led/controller.go#L94), [(controller *Controller).Update](../../led/controller.go#L100), [(controller *Controller).populateFixtureData](../../led/controller.go#L136), [(universe *universe).shouldSendPacket](../../led/controller.go#L166), [(universe *universe).markSent](../../led/controller.go#L173), [createBlankPacket](../../led/controller.go#L180), [(controller *Controller).sendPacket](../../led/controller.go#L287)

### [led/fixture.go](../../led/fixture.go)

Hardcoded physical fixture mapping for the 2026 Hub LEDs.

Fonksiyon yok; tip, sabit veya üretilmiş enum tanımı.

### [led/mode.go](../../led/mode.go)

Contains the set of display modes for the hub LEDs.

[getModeNames](../../led/mode.go#L51), [IsValidMode](../../led/mode.go#L67), [colorForMode](../../led/mode.go#L72)

### [led/zone.go](../../led/zone.go)

Represents one alliance's independently controlled Hub LED zone.

[(zone *zone).updatePixels](../../led/zone.go#L40), [(zone *zone).updateSingleColorMode](../../led/zone.go#L68), [(zone *zone).updatePulseMode](../../led/zone.go#L75), [(zone *zone).updateStartupMode](../../led/zone.go#L84), [(zone *zone).fillSide](../../led/zone.go#L108), [(zone *zone).fillFixture](../../led/zone.go#L124), [fillOrder](../../led/zone.go#L141), [(zone *zone).updateAdvantageMode](../../led/zone.go#L155), [(zone *zone).sweepFixture](../../led/zone.go#L170), [(zone *zone).updateRainbowMode](../../led/zone.go#L191), [(zone *zone).updateSideTestMode](../../led/zone.go#L245)

### [main.go](../../main.go)

[main](../../main.go#L18)

### [model/alliance.go](../../model/alliance.go)

Model and datastore CRUD methods for a playoff alliance.

[(database *Database).CreateAlliance](../../model/alliance.go#L23), [(database *Database).GetAllianceById](../../model/alliance.go#L27), [(database *Database).UpdateAlliance](../../model/alliance.go#L31), [(database *Database).DeleteAlliance](../../model/alliance.go#L35), [(database *Database).TruncateAlliances](../../model/alliance.go#L39), [(database *Database).GetAllAlliances](../../model/alliance.go#L43), [(database *Database).UpdateAllianceFromMatch](../../model/alliance.go#L58), [(database *Database).GetOffFieldTeamIds](../../model/alliance.go#L93), [(database *Database).getOffFieldTeamIdsForAlliance](../../model/alliance.go#L111)

### [model/award.go](../../model/award.go)

Model and datastore CRUD methods for an award.

[(award Award).TbaTypeEnumString](../../model/award.go#L22), [(database *Database).CreateAward](../../model/award.go#L37), [(database *Database).GetAwardById](../../model/award.go#L41), [(database *Database).UpdateAward](../../model/award.go#L45), [(database *Database).DeleteAward](../../model/award.go#L49), [(database *Database).TruncateAwards](../../model/award.go#L53), [(database *Database).GetAllAwards](../../model/award.go#L57), [(database *Database).GetAwardsByType](../../model/award.go#L71)

### [model/database.go](../../model/database.go)

Functions for manipulating the per-event Bolt datastore.

[OpenDatabase](../../model/database.go#L42), [(database *Database).Close](../../model/database.go#L94), [(database *Database).Backup](../../model/database.go#L99), [(database *Database).WriteBackup](../../model/database.go#L127)

### [model/event_settings.go](../../model/event_settings.go)

Model and datastore read/write methods for event-level configuration.

[(database *Database).GetEventSettings](../../model/event_settings.go#L124), [(database *Database).UpdateEventSettings](../../model/event_settings.go#L164)

### [model/judging_slot.go](../../model/judging_slot.go)

Model and datastore CRUD methods for a judging slot.

[(database *Database).CreateJudgingSlot](../../model/judging_slot.go#L24), [(database *Database).TruncateJudgingSlots](../../model/judging_slot.go#L28), [(database *Database).GetAllJudgingSlots](../../model/judging_slot.go#L32)

### [model/lower_third.go](../../model/lower_third.go)

Model and datastore CRUD methods for the text on a lower third slide.

[(database *Database).CreateLowerThird](../../model/lower_third.go#L20), [(database *Database).GetLowerThirdById](../../model/lower_third.go#L24), [(database *Database).UpdateLowerThird](../../model/lower_third.go#L28), [(database *Database).DeleteLowerThird](../../model/lower_third.go#L32), [(database *Database).TruncateLowerThirds](../../model/lower_third.go#L36), [(database *Database).GetAllLowerThirds](../../model/lower_third.go#L40), [(database *Database).GetLowerThirdsByAwardId](../../model/lower_third.go#L54), [(database *Database).GetNextLowerThirdDisplayOrder](../../model/lower_third.go#L69)

### [model/match.go](../../model/match.go)

Model and datastore CRUD methods for a match at an event.

[(t MatchType).Get](../../model/match.go#L27), [(database *Database).CreateMatch](../../model/match.go#L68), [(database *Database).GetMatchById](../../model/match.go#L72), [(database *Database).UpdateMatch](../../model/match.go#L76), [(database *Database).DeleteMatch](../../model/match.go#L80), [(database *Database).TruncateMatches](../../model/match.go#L84), [(database *Database).GetMatchByTypeOrder](../../model/match.go#L88), [(database *Database).GetMatchesByType](../../model/match.go#L102), [(match *Match).IsComplete](../../model/match.go#L124), [(match *Match).ShouldAllowSubstitution](../../model/match.go#L129), [(match *Match).ShouldAllowNexusSubstitution](../../model/match.go#L134), [(match *Match).ShouldUpdateCards](../../model/match.go#L139), [(match *Match).ShouldUpdateRankings](../../model/match.go#L144), [(match *Match).ShouldUpdatePlayoffMatches](../../model/match.go#L149), [(match *Match).IsLineupEqual](../../model/match.go#L154), [MatchTypeFromString](../../model/match.go#L159), [(key TbaMatchKey).String](../../model/match.go#L174)

### [model/match_result.go](../../model/match_result.go)

Model and datastore CRUD methods for the results (score and fouls) from a match at an event.

[NewMatchResult](../../model/match_result.go#L26), [(database *Database).CreateMatchResult](../../model/match_result.go#L35), [(database *Database).GetMatchResultForMatch](../../model/match_result.go#L39), [(database *Database).UpdateMatchResult](../../model/match_result.go#L55), [(database *Database).DeleteMatchResult](../../model/match_result.go#L59), [(database *Database).TruncateMatchResults](../../model/match_result.go#L63), [(matchResult *MatchResult).RedScoreSummary](../../model/match_result.go#L68), [(matchResult *MatchResult).BlueScoreSummary](../../model/match_result.go#L73), [(matchResult *MatchResult).CorrectPlayoffScore](../../model/match_result.go#L78)

### [model/matchtype_string.go](../../model/matchtype_string.go)

Code generated by "stringer -type=MatchType"; DO NOT EDIT.

[_](../../model/matchtype_string.go#L7), [(i MatchType).String](../../model/matchtype_string.go#L21)

### [model/ranking.go](../../model/ranking.go)

Model and datastore CRUD methods for team ranking data at an event.

[(database *Database).CreateRanking](../../model/ranking.go#L13), [(database *Database).GetRankingForTeam](../../model/ranking.go#L17), [(database *Database).UpdateRanking](../../model/ranking.go#L21), [(database *Database).DeleteRanking](../../model/ranking.go#L25), [(database *Database).TruncateRankings](../../model/ranking.go#L29), [(database *Database).GetAllRankings](../../model/ranking.go#L33), [(database *Database).ReplaceAllRankings](../../model/ranking.go#L48)

### [model/schedule_block.go](../../model/schedule_block.go)

Model and datastore CRUD methods for a schedule block at an event.

[(database *Database).CreateScheduleBlock](../../model/schedule_block.go#L21), [(database *Database).GetScheduleBlocksByMatchType](../../model/schedule_block.go#L25), [(database *Database).DeleteScheduleBlocksByMatchType](../../model/schedule_block.go#L47), [(database *Database).TruncateScheduleBlocks](../../model/schedule_block.go#L61)

### [model/scheduled_break.go](../../model/scheduled_break.go)

Model and datastore CRUD methods for a scheduled break at an event.

[(database *Database).CreateScheduledBreak](../../model/scheduled_break.go#L22), [(database *Database).GetScheduledBreakById](../../model/scheduled_break.go#L26), [(database *Database).UpdateScheduledBreak](../../model/scheduled_break.go#L30), [(database *Database).GetScheduledBreaksByMatchType](../../model/scheduled_break.go#L34), [(database *Database).GetScheduledBreakByMatchTypeOrder](../../model/scheduled_break.go#L56), [(database *Database).DeleteScheduledBreaksByMatchType](../../model/scheduled_break.go#L73), [(database *Database).TruncateScheduledBreaks](../../model/scheduled_break.go#L87)

### [model/sponsor_slide.go](../../model/sponsor_slide.go)

Model and datastore CRUD methods for the sponsor slideshow.

[(database *Database).CreateSponsorSlide](../../model/sponsor_slide.go#L20), [(database *Database).GetSponsorSlideById](../../model/sponsor_slide.go#L24), [(database *Database).UpdateSponsorSlide](../../model/sponsor_slide.go#L28), [(database *Database).DeleteSponsorSlide](../../model/sponsor_slide.go#L32), [(database *Database).TruncateSponsorSlides](../../model/sponsor_slide.go#L36), [(database *Database).GetAllSponsorSlides](../../model/sponsor_slide.go#L40), [(database *Database).GetNextSponsorSlideDisplayOrder](../../model/sponsor_slide.go#L54)

### [model/table.go](../../model/table.go)

Defines a "table" wrapper struct and helper methods for persisting data using Bolt.

[newTable](../../model/table.go#L28), [(table *table[R]).getById](../../model/table.go#L82), [(table *table[R]).getAll](../../model/table.go#L104), [(table *table[R]).create](../../model/table.go#L130), [(table *table[R]).update](../../model/table.go#L178), [(table *table[R]).delete](../../model/table.go#L210), [(table *table[R]).truncate](../../model/table.go#L231), [(table *table[R]).getBucket](../../model/table.go#L251), [idToKey](../../model/table.go#L260)

### [model/team.go](../../model/team.go)

Model and datastore CRUD methods for a team at an event.

[(database *Database).CreateTeam](../../model/team.go#L27), [(database *Database).GetTeamById](../../model/team.go#L31), [(database *Database).UpdateTeam](../../model/team.go#L35), [(database *Database).DeleteTeam](../../model/team.go#L39), [(database *Database).TruncateTeams](../../model/team.go#L43), [(database *Database).GetAllTeams](../../model/team.go#L47)

### [model/test_helpers.go](../../model/test_helpers.go)

Helper methods for use in tests in this package and others.

[SetupTestDb](../../model/test_helpers.go#L15), [BuildTestMatchResult](../../model/test_helpers.go#L29), [BuildTestAlliances](../../model/test_helpers.go#L38)

### [model/user_session.go](../../model/user_session.go)

Model and datastore CRUD methods for a user login session.

[(database *Database).CreateUserSession](../../model/user_session.go#L17), [(database *Database).GetUserSessionByToken](../../model/user_session.go#L21), [(database *Database).DeleteUserSession](../../model/user_session.go#L35), [(database *Database).TruncateUserSessions](../../model/user_session.go#L39)

### [network/access_point.go](../../network/access_point.go)

Methods for configuring a Vivid-Hosting VH-113 access point running OpenWRT for team SSIDs and VLANs.

[(ap *AccessPoint).SetSettings](../../network/access_point.go#L80), [(ap *AccessPoint).Run](../../network/access_point.go#L95), [(ap *AccessPoint).ConfigureTeamWifi](../../network/access_point.go#L114), [(ap *AccessPoint).updateMonitoring](../../network/access_point.go#L169), [(ap *AccessPoint).statusMatchesLastConfiguration](../../network/access_point.go#L229), [addStation](../../network/access_point.go#L247), [updateTeamWifiStatus](../../network/access_point.go#L258), [(apStatus *accessPointStatus).toLogString](../../network/access_point.go#L288)

### [network/sccswitch.go](../../network/sccswitch.go)

Methods for configuring an SCC Switch via SSH.

[NewSCCSwitch](../../network/sccswitch.go#L38), [(scc *SCCSwitch).SetTeamEthernetEnabled](../../network/sccswitch.go#L52), [(scc *SCCSwitch).runCommandSequence](../../network/sccswitch.go#L80)

### [network/switch.go](../../network/switch.go)

Methods for configuring a Cisco Switch 3500-series switch for team VLANs.

[NewSwitch](../../network/switch.go#L49), [(sw *Switch).ConfigureTeamEthernet](../../network/switch.go#L61), [(sw *Switch).runCommand](../../network/switch.go#L132), [(sw *Switch).runConfigCommand](../../network/switch.go#L167)

### [partner/blackmagic.go](../../partner/blackmagic.go)

Client for interfacing with one or more Blackmagic HyperDeck devices to automatically record matches.

[NewBlackmagicClient](../../partner/blackmagic.go#L27), [(client *BlackmagicClient).StartRecording](../../partner/blackmagic.go#L39), [(client *BlackmagicClient).StopRecording](../../partner/blackmagic.go#L44), [(client *BlackmagicClient).sendCommand](../../partner/blackmagic.go#L50)

### [partner/companion.go](../../partner/companion.go)

Client for interfacing with Bitfocus Companion to automatically trigger events during matches.

[NewCompanionClient](../../partner/companion.go#L48), [(client *CompanionClient).SendEvent](../../partner/companion.go#L61), [(client *CompanionClient).sendCommand](../../partner/companion.go#L77), [(client *CompanionClient).IsEnabled](../../partner/companion.go#L101), [(client *CompanionClient).GetEventConfig](../../partner/companion.go#L106)

### [partner/nexus.go](../../partner/nexus.go)

Methods for interfacing with the Nexus for FRC API.

[NewNexusClient](../../partner/nexus.go#L65), [(client *NexusClient).GetLineup](../../partner/nexus.go#L70), [(client *NexusClient).AutoQueue](../../partner/nexus.go#L121), [(client *NexusClient).MatchStarted](../../partner/nexus.go#L137), [(client *NexusClient).BreakStarted](../../partner/nexus.go#L143), [(client *NexusClient).BreakEnded](../../partner/nexus.go#L149), [(client *NexusClient).postAutoQueueRequest](../../partner/nexus.go#L155), [(client *NexusClient).getRequest](../../partner/nexus.go#L200)

### [partner/tba.go](../../partner/tba.go)

Methods for publishing data to and retrieving data from The Blue Alliance.

[stringPtr](../../partner/tba.go#L162), [NewTbaClient](../../partner/tba.go#L174), [(client *TbaClient).GetTeam](../../partner/tba.go#L184), [(client *TbaClient).GetRobotName](../../partner/tba.go#L204), [(client *TbaClient).GetTeamAwards](../../partner/tba.go#L231), [(client *TbaClient).DownloadTeamAvatar](../../partner/tba.go#L264), [(client *TbaClient).PublishTeams](../../partner/tba.go#L305), [(client *TbaClient).PublishMatches](../../partner/tba.go#L329), [(client *TbaClient).PublishRankings](../../partner/tba.go#L445), [(client *TbaClient).PublishAlliances](../../partner/tba.go#L482), [(client *TbaClient).PublishAwards](../../partner/tba.go#L525), [(client *TbaClient).DeletePublishedMatches](../../partner/tba.go#L556), [checkTbaPostResponse](../../partner/tba.go#L564), [(client *TbaClient).getEventName](../../partner/tba.go#L583), [getTbaTeam](../../partner/tba.go#L607), [(client *TbaClient).getRequest](../../partner/tba.go#L612), [(client *TbaClient).postRequest](../../partner/tba.go#L626), [createTbaAlliance](../../partner/tba.go#L657), [createTbaScoringBreakdown](../../partner/tba.go#L678)

### [playoff/alliance_source.go](../../playoff/alliance_source.go)

Represents how the alliance is determined to fill a given spot at a given stage in a playoff tournament.

[(source allianceSelectionSource).AllianceId](../../playoff/alliance_source.go#L33), [(source allianceSelectionSource).displayName](../../playoff/alliance_source.go#L37), [(source allianceSelectionSource).setDestination](../../playoff/alliance_source.go#L41), [(source allianceSelectionSource).update](../../playoff/alliance_source.go#L45), [(source allianceSelectionSource).traverse](../../playoff/alliance_source.go#L49), [(source matchupSource).AllianceId](../../playoff/alliance_source.go#L60), [(source matchupSource).displayName](../../playoff/alliance_source.go#L68), [(source matchupSource).setDestination](../../playoff/alliance_source.go#L75), [(source matchupSource).update](../../playoff/alliance_source.go#L86), [(source matchupSource).traverse](../../playoff/alliance_source.go#L93)

### [playoff/break_spec.go](../../playoff/break_spec.go)

Represents a scheduled break in the playoff match schedule.

Fonksiyon yok; tip, sabit veya üretilmiş enum tanımı.

### [playoff/double_elimination.go](../../playoff/double_elimination.go)

Defines the tournament structure for a double-elimination bracket culminating in a best-of-three final.

[newDoubleEliminationBracket](../../playoff/double_elimination.go#L15), [newFourAllianceDoubleEliminationBracket](../../playoff/double_elimination.go#L27), [newEightAllianceDoubleEliminationBracket](../../playoff/double_elimination.go#L91), [newDoubleEliminationMatch](../../playoff/double_elimination.go#L216)

### [playoff/match_group.go](../../playoff/match_group.go)

Interface representing a stage of a playoff tournament containing one or more related matches.

[collectMatchGroups](../../playoff/match_group.go#L44), [collectMatchSpecs](../../playoff/match_group.go#L59)

### [playoff/matchup.go](../../playoff/matchup.go)

Models and logic encapsulating a group of one or more matches between the same two alliances at a given point in a playoff tournament.

[(matchup *Matchup).Id](../../playoff/matchup.go#L30), [(matchup *Matchup).MatchSpecs](../../playoff/matchup.go#L34), [(matchup *Matchup).update](../../playoff/matchup.go#L38), [(matchup *Matchup).setSourceDestinations](../../playoff/matchup.go#L90), [(matchup *Matchup).traverse](../../playoff/matchup.go#L95), [(matchup *Matchup).RedAllianceSourceDisplayName](../../playoff/matchup.go#L110), [(matchup *Matchup).BlueAllianceSourceDisplayName](../../playoff/matchup.go#L116), [(matchup *Matchup).RedAllianceDestination](../../playoff/matchup.go#L121), [(matchup *Matchup).BlueAllianceDestination](../../playoff/matchup.go#L126), [(matchup *Matchup).StatusText](../../playoff/matchup.go#L131), [(matchup *Matchup).WinningAllianceId](../../playoff/matchup.go#L156), [(matchup *Matchup).LosingAllianceId](../../playoff/matchup.go#L167), [(matchup *Matchup).IsComplete](../../playoff/matchup.go#L178), [(matchup *Matchup).IsLosingAllianceEliminated](../../playoff/matchup.go#L183), [(matchup *Matchup).isFinal](../../playoff/matchup.go#L188), [(matchup *Matchup).allianceDestination](../../playoff/matchup.go#L193), [formatDestinationMatchName](../../playoff/matchup.go#L217)

### [playoff/playoff_match_result.go](../../playoff/playoff_match_result.go)

Represents the outcome of a given playoff match.

Fonksiyon yok; tip, sabit veya üretilmiş enum tanımı.

### [playoff/playoff_tournament.go](../../playoff/playoff_tournament.go)

Models and logic encapsulating the common aspects of all supported playoff tournament formats.

[NewPlayoffTournament](../../playoff/playoff_tournament.go#L24), [(tournament *PlayoffTournament).MatchGroups](../../playoff/playoff_tournament.go#L64), [(tournament *PlayoffTournament).FinalMatchup](../../playoff/playoff_tournament.go#L69), [(tournament *PlayoffTournament).IsComplete](../../playoff/playoff_tournament.go#L74), [(tournament *PlayoffTournament).WinningAllianceId](../../playoff/playoff_tournament.go#L80), [(tournament *PlayoffTournament).FinalistAllianceId](../../playoff/playoff_tournament.go#L86), [(tournament *PlayoffTournament).Traverse](../../playoff/playoff_tournament.go#L91), [(tournament *PlayoffTournament).CreateMatchesAndBreaks](../../playoff/playoff_tournament.go#L97), [(tournament *PlayoffTournament).UpdateMatches](../../playoff/playoff_tournament.go#L186), [positionRedTeams](../../playoff/playoff_tournament.go#L254), [positionBlueTeams](../../playoff/playoff_tournament.go#L261)

### [playoff/single_elimination.go](../../playoff/single_elimination.go)

Defines the tournament structure for a single-elimination, best-of-three bracket.

[newSingleEliminationBracket](../../playoff/single_elimination.go#L16), [newSingleEliminationAllianceSource](../../playoff/single_elimination.go#L206), [newSingleEliminationMatch](../../playoff/single_elimination.go#L220), [newFinalMatches](../../playoff/single_elimination.go#L232)

### [playoff/test_helpers.go](../../playoff/test_helpers.go)

Helper methods for use in tests in this package and others.

[setupTestDb](../../playoff/test_helpers.go#L14), [assertMatchSpecs](../../playoff/test_helpers.go#L31), [assertMatchSpecAlliances](../../playoff/test_helpers.go#L59), [assertMatchGroups](../../playoff/test_helpers.go#L72), [assertMatch](../../playoff/test_helpers.go#L83), [assertBreak](../../playoff/test_helpers.go#L132), [assertMatchupOutcome](../../playoff/test_helpers.go#L147)

### [plc/armorblock_string.go](../../plc/armorblock_string.go)

Code generated by "stringer -type=armorBlock"; DO NOT EDIT.

[_](../../plc/armorblock_string.go#L7), [(i armorBlock).String](../../plc/armorblock_string.go#L22)

### [plc/coil_string.go](../../plc/coil_string.go)

Code generated by "stringer -type=coil"; DO NOT EDIT.

[_](../../plc/coil_string.go#L7), [(i coil).String](../../plc/coil_string.go#L31)

### [plc/input_string.go](../../plc/input_string.go)

Code generated by "stringer -type=input"; DO NOT EDIT.

[_](../../plc/input_string.go#L7), [(i input).String](../../plc/input_string.go#L46)

### [plc/plc.go](../../plc/plc.go)

Methods for interfacing with the field PLC.

[(plc *ModbusPlc).SetAddress](../../plc/plc.go#L162), [(plc *ModbusPlc).IsEnabled](../../plc/plc.go#L173), [(plc *ModbusPlc).IsHealthy](../../plc/plc.go#L178), [(plc *ModbusPlc).IoChangeNotifier](../../plc/plc.go#L183), [(plc *ModbusPlc).Run](../../plc/plc.go#L188), [(plc *ModbusPlc).GetArmorBlockStatuses](../../plc/plc.go#L212), [(plc *ModbusPlc).GetFieldEStop](../../plc/plc.go#L221), [(plc *ModbusPlc).GetTeamEStops](../../plc/plc.go#L226), [(plc *ModbusPlc).GetTeamAStops](../../plc/plc.go#L238), [(plc *ModbusPlc).GetEthernetConnected](../../plc/plc.go#L250), [(plc *ModbusPlc).IsFtaReady](../../plc/plc.go#L264), [(plc *ModbusPlc).ResetMatch](../../plc/plc.go#L269), [(plc *ModbusPlc).SetStackLights](../../plc/plc.go#L282), [(plc *ModbusPlc).SetStackBuzzer](../../plc/plc.go#L290), [(plc *ModbusPlc).SetFieldResetLight](../../plc/plc.go#L295), [(plc *ModbusPlc).SetAwardsModeLight](../../plc/plc.go#L300), [(plc *ModbusPlc).GetHubCounts](../../plc/plc.go#L305), [(plc *ModbusPlc).SetHubMotors](../../plc/plc.go#L310), [(plc *ModbusPlc).SetHubLights](../../plc/plc.go#L316), [(plc *ModbusPlc).GetCycleState](../../plc/plc.go#L321), [(plc *ModbusPlc).GetInputNames](../../plc/plc.go#L325), [(plc *ModbusPlc).GetRegisterNames](../../plc/plc.go#L333), [(plc *ModbusPlc).GetCoilNames](../../plc/plc.go#L341), [(plc *ModbusPlc).SetCoilOverride](../../plc/plc.go#L349), [(plc *ModbusPlc).ClearCoilOverride](../../plc/plc.go#L359), [(plc *ModbusPlc).connect](../../plc/plc.go#L366), [(plc *ModbusPlc).resetConnection](../../plc/plc.go#L383), [(plc *ModbusPlc).update](../../plc/plc.go#L393), [(plc *ModbusPlc).readInputs](../../plc/plc.go#L423), [(plc *ModbusPlc).readRegisters](../../plc/plc.go#L442), [(plc *ModbusPlc).writeCoils](../../plc/plc.go#L465), [(plc *ModbusPlc).generateIoChangeMessage](../../plc/plc.go#L485), [(plc *ModbusPlc).getEffectiveCoils](../../plc/plc.go#L496), [(plc *ModbusPlc).getCoilOverrideStates](../../plc/plc.go#L506), [byteToBool](../../plc/plc.go#L521), [byteToUint](../../plc/plc.go#L532), [boolToByte](../../plc/plc.go#L540)

### [plc/register_string.go](../../plc/register_string.go)

Code generated by "stringer -type=register"; DO NOT EDIT.

[_](../../plc/register_string.go#L7), [(i register).String](../../plc/register_string.go#L29)

### [tournament/alliance_cards.go](../../tournament/alliance_cards.go)

[CalculateAllianceCards](../../tournament/alliance_cards.go#L11)

### [tournament/awards.go](../../tournament/awards.go)

Functions for managing awards and their associated lower thirds.

[CreateOrUpdateAward](../../tournament/awards.go#L14), [DeleteAward](../../tournament/awards.go#L80), [CreateOrUpdateWinnerAndFinalistAwards](../../tournament/awards.go#L100), [createOrUpdateAwardLowerThird](../../tournament/awards.go#L168)

### [tournament/judging_schedule.go](../../tournament/judging_schedule.go)

Logic for creating judging schedules.

[BuildJudgingSchedule](../../tournament/judging_schedule.go#L43), [createTeamMatchMap](../../tournament/judging_schedule.go#L201), [getNextSlotForTeam](../../tournament/judging_schedule.go#L220)

### [tournament/qualification_rankings.go](../../tournament/qualification_rankings.go)

Functions for calculating the qualification rankings.

[CalculateRankings](../../tournament/qualification_rankings.go#L17), [CalculateTeamCards](../../tournament/qualification_rankings.go#L81), [addMatchResultToRankings](../../tournament/qualification_rankings.go#L135), [sortRankings](../../tournament/qualification_rankings.go#L161)

### [tournament/schedule.go](../../tournament/schedule.go)

Functions for creating practice and qualification match schedules.

[BuildRandomSchedule](../../tournament/schedule.go#L28), [countMatches](../../tournament/schedule.go#L111)

### [tournament/test_helpers.go](../../tournament/test_helpers.go)

Helper methods for use in tests in this package and others.

[CreateTestAlliances](../../tournament/test_helpers.go#L13), [setupTestDb](../../tournament/test_helpers.go#L24)

### [web/alliance_selection.go](../../web/alliance_selection.go)

Web routes for conducting the alliance selection process.

[(web *Web).allianceSelectionGetHandler](../../web/alliance_selection.go#L33), [(web *Web).allianceSelectionPostHandler](../../web/alliance_selection.go#L42), [(web *Web).allianceSelectionStartHandler](../../web/alliance_selection.go#L103), [(web *Web).allianceSelectionResetHandler](../../web/alliance_selection.go#L148), [(web *Web).allianceSelectionFinalizeHandler](../../web/alliance_selection.go#L178), [(web *Web).allianceSelectionWebsocketHandler](../../web/alliance_selection.go#L274), [(web *Web).renderAllianceSelection](../../web/alliance_selection.go#L368), [(web *Web).canModifyAllianceSelection](../../web/alliance_selection.go#L413), [(web *Web).canResetAllianceSelection](../../web/alliance_selection.go#L422), [(web *Web).determineNextCell](../../web/alliance_selection.go#L436)

### [web/alliance_station_display.go](../../web/alliance_station_display.go)

Web handlers for the alliance station display.

[(web *Web).allianceStationDisplayHandler](../../web/alliance_station_display.go#L15), [(web *Web).allianceStationDisplayWebsocketHandler](../../web/alliance_station_display.go#L37)

### [web/announcer_display.go](../../web/announcer_display.go)

Web handlers for announcer display.

[(web *Web).announcerDisplayHandler](../../web/announcer_display.go#L15), [(web *Web).announcerDisplayMatchLoadHandler](../../web/announcer_display.go#L37), [(web *Web).announcerDisplayScorePostedHandler](../../web/announcer_display.go#L52), [(web *Web).announcerDisplayWebsocketHandler](../../web/announcer_display.go#L67)

### [web/api.go](../../web/api.go)

Web API for providing JSON-formatted event data.

[(web *Web).matchesApiHandler](../../web/api.go#L51), [(web *Web).sponsorSlidesApiHandler](../../web/api.go#L97), [(web *Web).rankingsApiHandler](../../web/api.go#L124), [(web *Web).alliancesApiHandler](../../web/api.go#L185), [(web *Web).arenaWebsocketApiHandler](../../web/api.go#L208), [(web *Web).teamAvatarsApiHandler](../../web/api.go#L221), [(web *Web).bracketSvgApiHandler](../../web/api.go#L236), [(web *Web).generateBracketSvg](../../web/api.go#L254)

### [web/audience_display.go](../../web/audience_display.go)

Web handlers for audience screen display.

[(web *Web).audienceDisplayHandler](../../web/audience_display.go#L16), [(web *Web).audienceDisplayWebsocketHandler](../../web/audience_display.go#L48)

### [web/bracket_display.go](../../web/bracket_display.go)

Web handlers for the bracket display.

[(web *Web).bracketDisplayHandler](../../web/bracket_display.go#L15), [(web *Web).bracketDisplayWebsocketHandler](../../web/bracket_display.go#L36)

### [web/display_utils.go](../../web/display_utils.go)

Common utility methods for display web routes.

[(web *Web).enforceDisplayConfiguration](../../web/display_utils.go#L18), [(web *Web).registerDisplay](../../web/display_utils.go#L58)

### [web/field_monitor_display.go](../../web/field_monitor_display.go)

Web handlers for the field monitor display showing robot connection status.

[(web *Web).fieldMonitorDisplayHandler](../../web/field_monitor_display.go#L18), [(web *Web).fmsFieldMonitorDisplayHandler](../../web/field_monitor_display.go#L43), [(web *Web).fieldMonitorDisplayWebsocketHandler](../../web/field_monitor_display.go#L68)

### [web/login.go](../../web/login.go)

Web routes for authenticating with the server.

[(web *Web).loginHandler](../../web/login.go#L19), [(web *Web).loginPostHandler](../../web/login.go#L24), [(web *Web).renderLogin](../../web/login.go#L45), [(web *Web).userIsAdmin](../../web/login.go#L63), [(web *Web).getUserSessionFromCookie](../../web/login.go#L81), [(web *Web).checkAuthPassword](../../web/login.go#L93)

### [web/logo_display.go](../../web/logo_display.go)

Web routes for a display to show a static logo and configurable message.

[(web *Web).logoDisplayHandler](../../web/logo_display.go#L15), [(web *Web).logoDisplayWebsocketHandler](../../web/logo_display.go#L36)

### [web/match_logs.go](../../web/match_logs.go)

Web routes for viewing match logs

[(web *Web).matchLogsHandler](../../web/match_logs.go#L70), [(web *Web).matchLogsViewGetHandler](../../web/match_logs.go#L114), [(web *Web).getMatchLogFromRequest](../../web/match_logs.go#L144), [parseOptionalString](../../web/match_logs.go#L270), [parseOptionalBool](../../web/match_logs.go#L279), [parseOptionalFloat](../../web/match_logs.go#L292), [parseOptionalInt](../../web/match_logs.go#L305), [(web *Web).buildMatchLogsList](../../web/match_logs.go#L318)

### [web/match_play.go](../../web/match_play.go)

Web routes for controlling match play.

[(web *Web).matchPlayHandler](../../web/match_play.go#L36), [(web *Web).matchPlayMatchLoadHandler](../../web/match_play.go#L65), [(web *Web).matchPlayWebsocketHandler](../../web/match_play.go#L116), [(web *Web).commitPostAndLoadNextMatch](../../web/match_play.go#L378), [(web *Web).commitMatchScore](../../web/match_play.go#L402), [(web *Web).getCurrentMatchResult](../../web/match_play.go#L539), [(web *Web).commitCurrentMatchScore](../../web/match_play.go#L553), [(list MatchPlayList).Len](../../web/match_play.go#L558), [(list MatchPlayList).Less](../../web/match_play.go#L563), [(list MatchPlayList).Swap](../../web/match_play.go#L568), [(web *Web).buildMatchPlayList](../../web/match_play.go#L573)

### [web/match_review.go](../../web/match_review.go)

Web routes for editing match results.

[(web *Web).matchReviewHandler](../../web/match_review.go#L45), [(web *Web).matchReviewEditGetHandler](../../web/match_review.go#L89), [(web *Web).matchReviewSummaryPostHandler](../../web/match_review.go#L142), [(web *Web).matchReviewEditPostHandler](../../web/match_review.go#L183), [(web *Web).getMatchResultFromRequest](../../web/match_review.go#L229), [normalizeMatchResult](../../web/match_review.go#L260), [(web *Web).buildMatchReviewList](../../web/match_review.go#L285)

### [web/placeholder_display.go](../../web/placeholder_display.go)

Web routes for a placeholder display to be later configured by the server.

[(web *Web).placeholderDisplayHandler](../../web/placeholder_display.go#L15), [(web *Web).placeholderDisplayWebsocketHandler](../../web/placeholder_display.go#L36)

### [web/queueing_display.go](../../web/queueing_display.go)

Web handlers for queueing display.

[(web *Web).queueingDisplayHandler](../../web/queueing_display.go#L22), [(web *Web).queueingDisplayMatchLoadHandler](../../web/queueing_display.go#L46), [(web *Web).queueingDisplayWebsocketHandler](../../web/queueing_display.go#L110)

### [web/rankings_display.go](../../web/rankings_display.go)

Web handlers for the rankings display.

[(web *Web).rankingsDisplayHandler](../../web/rankings_display.go#L15), [(web *Web).rankingsDisplayWebsocketHandler](../../web/rankings_display.go#L36)

### [web/referee_panel.go](../../web/referee_panel.go)

Web handlers for the referee interface.

[(web *Web).refereePanelHandler](../../web/referee_panel.go#L23), [(web *Web).refereePanelFoulListHandler](../../web/referee_panel.go#L45), [(web *Web).refereePanelWebsocketHandler](../../web/referee_panel.go#L71)

### [web/report_pdf.go](../../web/report_pdf.go)

[newReportPdf](../../web/report_pdf.go#L35), [registerReportFont](../../web/report_pdf.go#L60), [(pdf *reportPdf).SetFont](../../web/report_pdf.go#L67), [findArialFontFiles](../../web/report_pdf.go#L87), [arialFontStyle](../../web/report_pdf.go#L126), [normalizeFontName](../../web/report_pdf.go#L149)

### [web/reports.go](../../web/reports.go)

Web handlers for generating CSV and PDF reports.

[(web *Web).rankingsCsvReportHandler](../../web/reports.go#L25), [(web *Web).rankingsPdfReportHandler](../../web/reports.go#L55), [(web *Web).findBackupTeams](../../web/reports.go#L128), [(web *Web).backupTeamsCsvReportHandler](../../web/reports.go#L173), [(web *Web).backupsPdfReportHandler](../../web/reports.go#L223), [(web *Web).couponsPdfReportHandler](../../web/reports.go#L290), [drawCoupon](../../web/reports.go#L324), [drawEventWatermark](../../web/reports.go#L336), [drawCenteredText](../../web/reports.go#L354), [drawPdfLogo](../../web/reports.go#L359), [(web *Web).scheduleCsvReportHandler](../../web/reports.go#L374), [(web *Web).schedulePdfReportHandler](../../web/reports.go#L410), [(web *Web).teamsCsvReportHandler](../../web/reports.go#L554), [(web *Web).teamsPdfReportHandler](../../web/reports.go#L584), [(web *Web).wpaKeysCsvReportHandler](../../web/reports.go#L661), [(web *Web).alliancesPdfReportHandler](../../web/reports.go#L684), [(web *Web).bracketPdfReportHandler](../../web/reports.go#L807), [surrogateText](../../web/reports.go#L836), [(web *Web).cyclePdfReportHandler](../../web/reports.go#L845), [(web *Web).ftaCsvReportHandler](../../web/reports.go#L949), [(web *Web).judgingSchedulePdfReportHandler](../../web/reports.go#L979), [addTimeGeneratedFooter](../../web/reports.go#L1102), [drawMultiLineCell](../../web/reports.go#L1111)

### [web/scoring_panel.go](../../web/scoring_panel.go)

Web handlers for scoring interface.

[(web *Web).scoringPanelHandler](../../web/scoring_panel.go#L37), [(web *Web).scoringPanelWebsocketHandler](../../web/scoring_panel.go#L67)

### [web/setup_awards.go](../../web/setup_awards.go)

Web routes for managing awards.

[(web *Web).awardsGetHandler](../../web/setup_awards.go#L17), [(web *Web).awardsPostHandler](../../web/setup_awards.go#L54)

### [web/setup_breaks.go](../../web/setup_breaks.go)

Web routes for managing scheduled breaks.

[(web *Web).breaksGetHandler](../../web/setup_breaks.go#L15), [(web *Web).breaksPostHandler](../../web/setup_breaks.go#L42)

### [web/setup_displays.go](../../web/setup_displays.go)

Web routes for configuring the field displays.

[(web *Web).displaysGetHandler](../../web/setup_displays.go#L20), [(web *Web).displaysWebsocketHandler](../../web/setup_displays.go#L42)

### [web/setup_field_testing.go](../../web/setup_field_testing.go)

Web routes for testing the field sounds, LEDs, and PLC.

[(web *Web).fieldTestingGetHandler](../../web/setup_field_testing.go#L28), [(web *Web).fieldTestingWebsocketHandler](../../web/setup_field_testing.go#L69), [fieldTestingOverridesAllowed](../../web/setup_field_testing.go#L187)

### [web/setup_judging.go](../../web/setup_judging.go)

Web routes for generating judging schedules.

[(web *Web).judgingGetHandler](../../web/setup_judging.go#L25), [(web *Web).judgingGeneratePostHandler](../../web/setup_judging.go#L34), [(web *Web).judgingClearPostHandler](../../web/setup_judging.go#L92), [(web *Web).renderJudging](../../web/setup_judging.go#L106)

### [web/setup_lower_thirds.go](../../web/setup_lower_thirds.go)

Web routes for managing lower thirds.

[(web *Web).lowerThirdsGetHandler](../../web/setup_lower_thirds.go#L19), [(web *Web).lowerThirdsWebsocketHandler](../../web/setup_lower_thirds.go#L48), [(web *Web).saveLowerThird](../../web/setup_lower_thirds.go#L158), [(web *Web).reorderLowerThird](../../web/setup_lower_thirds.go#L180)

### [web/setup_schedule.go](../../web/setup_schedule.go)

Web routes for generating practice and qualification schedules.

[(web *Web).scheduleGetHandler](../../web/setup_schedule.go#L22), [(web *Web).scheduleGeneratePostHandler](../../web/setup_schedule.go#L38), [(web *Web).scheduleSavePostHandler](../../web/setup_schedule.go#L122), [(web *Web).renderSchedule](../../web/setup_schedule.go#L171), [getScheduleBlocks](../../web/setup_schedule.go#L220), [getMatchType](../../web/setup_schedule.go#L247)

### [web/setup_settings.go](../../web/setup_settings.go)

Web routes for configuring the event settings.

[(web *Web).settingsGetHandler](../../web/setup_settings.go#L23), [(web *Web).settingsPostHandler](../../web/setup_settings.go#L32), [settingsSaveAllowed](../../web/setup_settings.go#L213), [settingsTabFromRequest](../../web/setup_settings.go#L217), [(web *Web).saveDbHandler](../../web/setup_settings.go#L227), [(web *Web).restoreDbHandler](../../web/setup_settings.go#L244), [(web *Web).clearDbHandler](../../web/setup_settings.go#L330), [(web *Web).settingsPublishAlliancesHandler](../../web/setup_settings.go#L381), [(web *Web).settingsPublishAwardsHandler](../../web/setup_settings.go#L403), [(web *Web).settingsPublishMatchesHandler](../../web/setup_settings.go#L425), [(web *Web).settingsPublishRankingsHandler](../../web/setup_settings.go#L454), [(web *Web).settingsPublishTeamsHandler](../../web/setup_settings.go#L476), [(web *Web).renderSettings](../../web/setup_settings.go#L497), [(web *Web).renderSettingsWithStatus](../../web/setup_settings.go#L501), [(web *Web).deleteMatchDataForType](../../web/setup_settings.go#L526)

### [web/setup_sponsor_slides.go](../../web/setup_sponsor_slides.go)

Web routes for managing sponsor slides.

[(web *Web).sponsorSlidesGetHandler](../../web/setup_sponsor_slides.go#L15), [(web *Web).sponsorSlidesPostHandler](../../web/setup_sponsor_slides.go#L46), [(web *Web).reorderSponsorSlide](../../web/setup_sponsor_slides.go#L104)

### [web/setup_teams.go](../../web/setup_teams.go)

Web routes for configuring the team list.

[(web *Web).teamsGetHandler](../../web/setup_teams.go#L26), [(web *Web).teamsPostHandler](../../web/setup_teams.go#L35), [(web *Web).teamsRefreshHandler](../../web/setup_teams.go#L76), [(web *Web).teamsClearHandler](../../web/setup_teams.go#L108), [(web *Web).teamEditGetHandler](../../web/setup_teams.go#L127), [(web *Web).teamEditPostHandler](../../web/setup_teams.go#L160), [(web *Web).teamDeletePostHandler](../../web/setup_teams.go#L202), [(web *Web).teamsGenerateWpaKeysHandler](../../web/setup_teams.go#L231), [(web *Web).teamsUpdateProgressBarHandler](../../web/setup_teams.go#L260), [(web *Web).renderTeams](../../web/setup_teams.go#L268), [(web *Web).canModifyTeamList](../../web/setup_teams.go#L293), [(web *Web).populateOfficialTeamInfo](../../web/setup_teams.go#L302)

### [web/twitch_display.go](../../web/twitch_display.go)

Web routes for a display to show a configurable Twitch live video stream.

[(web *Web).twitchDisplayHandler](../../web/twitch_display.go#L15), [(web *Web).twitchDisplayWebsocketHandler](../../web/twitch_display.go#L36)

### [web/unpicked_display.go](../../web/unpicked_display.go)

Web handlers for unpicked teams display.

[(web *Web).unpickedDisplayHandler](../../web/unpicked_display.go#L13), [(web *Web).unpickedDisplayWebsocketHandler](../../web/unpicked_display.go#L38)

### [web/wall_display.go](../../web/wall_display.go)

Web handlers for a display that can be shown on a screen supplemental to one with the audience display overlay.

[(web *Web).wallDisplayHandler](../../web/wall_display.go#L16), [(web *Web).wallDisplayWebsocketHandler](../../web/wall_display.go#L49)

### [web/web.go](../../web/web.go)

Configuration and functions for the event server web interface.

[NewWeb](../../web/web.go#L31), [(web *Web).ServeWebInterface](../../web/web.go#L89), [(web *Web).indexHandler](../../web/web.go#L101), [addNoCacheHeader](../../web/web.go#L118), [(web *Web).newHandler](../../web/web.go#L128), [handleWebErr](../../web/web.go#L249), [(web *Web).parseFiles](../../web/web.go#L255)

### [web/webpage_display.go](../../web/webpage_display.go)

Web routes for a display to show an arbitrary web page.

[(web *Web).webpageDisplayHandler](../../web/webpage_display.go#L15), [(web *Web).webpageDisplayWebsocketHandler](../../web/webpage_display.go#L36)

### [web/websocket_helpers.go](../../web/websocket_helpers.go)

[writeWebsocketError](../../web/websocket_helpers.go#L11), [writeWebsocketMessage](../../web/websocket_helpers.go#L17), [closeWebsocket](../../web/websocket_helpers.go#L23)

### [web/youtube_display.go](../../web/youtube_display.go)

Web routes for a display to show a configurable YouTube live video stream.

[(web *Web).youtubeDisplayHandler](../../web/youtube_display.go#L15), [(web *Web).youtubeDisplayWebsocketHandler](../../web/youtube_display.go#L36)

### [websocket/notifier.go](../../websocket/notifier.go)

Publish-subscribe model for nonblocking notification of server events to websocket clients.

[NewNotifier](../../websocket/notifier.go#L28), [(notifier *Notifier).Notify](../../websocket/notifier.go#L36), [(notifier *Notifier).NotifyWithMessage](../../websocket/notifier.go#L42), [(notifier *Notifier).notifyListener](../../websocket/notifier.go#L52), [(notifier *Notifier).listen](../../websocket/notifier.go#L72), [(notifier *Notifier).getMessageBody](../../websocket/notifier.go#L82)

### [websocket/websocket.go](../../websocket/websocket.go)

Functions for the server side of handling websockets.

[NewWebsocket](../../websocket/websocket.go#L37), [NewTestWebsocket](../../websocket/websocket.go#L45), [(ws *Websocket).Close](../../websocket/websocket.go#L49), [(ws *Websocket).Read](../../websocket/websocket.go#L53), [(ws *Websocket).ReadWithTimeout](../../websocket/websocket.go#L72), [(ws *Websocket).Write](../../websocket/websocket.go#L92), [(ws *Websocket).WriteNotifier](../../websocket/websocket.go#L105), [(ws *Websocket).WriteError](../../websocket/websocket.go#L109), [(ws *Websocket).HandleNotifiers](../../websocket/websocket.go#L114)

## Handler → şablon → tarayıcı dosyaları

Tablo yalnızca handler gövdesinde doğrudan geçen şablonları listeler. `renderSettings`, `renderTeams` gibi yardımcıların şablonları sonraki tam şablon envanterinde görünür.

| Rota | Doğrudan kullanılan şablon |
| --- | --- |
| `GET /` | [templates/base.html](../../templates/base.html), [templates/index.html](../../templates/index.html) |
| `GET /display` | [templates/placeholder_display.html](../../templates/placeholder_display.html) |
| `GET /displays/alliance_station` | [templates/alliance_station_display.html](../../templates/alliance_station_display.html) |
| `GET /displays/announcer` | [templates/announcer_display.html](../../templates/announcer_display.html), [templates/base.html](../../templates/base.html) |
| `GET /displays/announcer/match_load` | [templates/announcer_display_match_load.html](../../templates/announcer_display_match_load.html) |
| `GET /displays/announcer/score_posted` | [templates/announcer_display_score_posted.html](../../templates/announcer_display_score_posted.html) |
| `GET /displays/audience` | [templates/audience_display.html](../../templates/audience_display.html) |
| `GET /displays/bracket` | [templates/bracket_display.html](../../templates/bracket_display.html) |
| `GET /displays/field_monitor` | [templates/field_monitor_display.html](../../templates/field_monitor_display.html) |
| `GET /displays/fms_field_monitor` | [templates/fms_field_monitor_display.html](../../templates/fms_field_monitor_display.html) |
| `GET /displays/logo` | [templates/logo_display.html](../../templates/logo_display.html) |
| `GET /displays/queueing` | [templates/queueing_display.html](../../templates/queueing_display.html) |
| `GET /displays/queueing/match_load` | [templates/queueing_display_match_load.html](../../templates/queueing_display_match_load.html) |
| `GET /displays/rankings` | [templates/rankings_display.html](../../templates/rankings_display.html) |
| `GET /displays/twitch` | [templates/twitch_display.html](../../templates/twitch_display.html) |
| `GET /displays/unpicked` | [templates/unpicked_display.html](../../templates/unpicked_display.html) |
| `GET /displays/wall` | [templates/wall_display.html](../../templates/wall_display.html) |
| `GET /displays/webpage` | [templates/webpage_display.html](../../templates/webpage_display.html) |
| `GET /displays/youtube` | [templates/youtube_display.html](../../templates/youtube_display.html) |
| `GET /match_play` | [templates/audience_display_radio_buttons.html](../../templates/audience_display_radio_buttons.html), [templates/base.html](../../templates/base.html), [templates/match_play.html](../../templates/match_play.html) |
| `GET /match_play/match_load` | [templates/match_play_match_load.html](../../templates/match_play_match_load.html) |
| `GET /match_logs` | [templates/base.html](../../templates/base.html), [templates/match_logs.html](../../templates/match_logs.html) |
| `GET /match_logs/{matchId}/{stationId}/log` | [templates/base.html](../../templates/base.html), [templates/view_match_log.html](../../templates/view_match_log.html) |
| `GET /match_review` | [templates/base.html](../../templates/base.html), [templates/match_review.html](../../templates/match_review.html) |
| `GET /match_review/{matchId}/edit` | [templates/base.html](../../templates/base.html), [templates/edit_match_result.html](../../templates/edit_match_result.html) |
| `GET /panels/scoring/{position}` | [templates/base.html](../../templates/base.html), [templates/scoring_panel.html](../../templates/scoring_panel.html) |
| `GET /panels/referee` | [templates/base.html](../../templates/base.html), [templates/referee_panel.html](../../templates/referee_panel.html) |
| `GET /panels/referee/foul_list` | [templates/referee_panel_foul_list.html](../../templates/referee_panel_foul_list.html) |
| `GET /reports/csv/backups` | [templates/backups.csv](../../templates/backups.csv) |
| `GET /reports/csv/fta` | [templates/fta.csv](../../templates/fta.csv) |
| `GET /reports/csv/rankings` | [templates/rankings.csv](../../templates/rankings.csv) |
| `GET /reports/csv/schedule/{type}` | [templates/schedule.csv](../../templates/schedule.csv) |
| `GET /reports/csv/teams` | [templates/teams.csv](../../templates/teams.csv) |
| `GET /reports/pdf/bracket` | [templates/bracket_report.html](../../templates/bracket_report.html) |
| `GET /setup/awards` | [templates/base.html](../../templates/base.html), [templates/setup_awards.html](../../templates/setup_awards.html) |
| `GET /setup/breaks` | [templates/base.html](../../templates/base.html), [templates/setup_breaks.html](../../templates/setup_breaks.html) |
| `GET /setup/displays` | [templates/base.html](../../templates/base.html), [templates/setup_displays.html](../../templates/setup_displays.html) |
| `GET /setup/field_testing` | [templates/base.html](../../templates/base.html), [templates/setup_field_testing.html](../../templates/setup_field_testing.html) |
| `GET /setup/lower_thirds` | [templates/audience_display_radio_buttons.html](../../templates/audience_display_radio_buttons.html), [templates/base.html](../../templates/base.html), [templates/setup_lower_thirds.html](../../templates/setup_lower_thirds.html) |
| `GET /setup/sponsor_slides` | [templates/base.html](../../templates/base.html), [templates/setup_sponsor_slides.html](../../templates/setup_sponsor_slides.html) |
| `GET /setup/teams/{id}/edit` | [templates/base.html](../../templates/base.html), [templates/edit_team.html](../../templates/edit_team.html) |

## Tüm şablonlar

| Şablon | Doğrudan JS/CSS referansları |
| --- | --- |
| [templates/alliance_selection.html](../../templates/alliance_selection.html) | [static/js/alliance_selection.js](../../static/js/alliance_selection.js), [static/js/match_timing.js](../../static/js/match_timing.js) |
| [templates/alliance_station_display.html](../../templates/alliance_station_display.html) | [static/css/alliance_station_display.css](../../static/css/alliance_station_display.css), [static/css/cheesy-arena.css](../../static/css/cheesy-arena.css), [static/css/lib/bootstrap.min.css](../../static/css/lib/bootstrap.min.css), [static/js/alliance_station_display.js](../../static/js/alliance_station_display.js), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/match_timing.js](../../static/js/match_timing.js) |
| [templates/announcer_display.html](../../templates/announcer_display.html) | [static/js/announcer_display.js](../../static/js/announcer_display.js), [static/js/match_timing.js](../../static/js/match_timing.js) |
| [templates/announcer_display_match_load.html](../../templates/announcer_display_match_load.html) |  |
| [templates/announcer_display_score_posted.html](../../templates/announcer_display_score_posted.html) |  |
| [templates/audience_display.html](../../templates/audience_display.html) | [static/css/audience_display.css](../../static/css/audience_display.css), [static/css/cheesy-arena.css](../../static/css/cheesy-arena.css), [static/css/lib/bootstrap.min.css](../../static/css/lib/bootstrap.min.css), [static/js/audience_display.js](../../static/js/audience_display.js), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/lib/bootstrap.bundle.min.js](../../static/js/lib/bootstrap.bundle.min.js), [static/js/lib/handlebars-1.3.0.js](../../static/js/lib/handlebars-1.3.0.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.transit.min.js](../../static/js/lib/jquery.transit.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/match_timing.js](../../static/js/match_timing.js) |
| [templates/audience_display_radio_buttons.html](../../templates/audience_display_radio_buttons.html) |  |
| [templates/backups.csv](../../templates/backups.csv) |  |
| [templates/base.html](../../templates/base.html) | [static/css/cheesy-arena.css](../../static/css/cheesy-arena.css), [static/css/lib/bootstrap-icons.min.css](../../static/css/lib/bootstrap-icons.min.css), [static/css/lib/bootstrap.min.css](../../static/css/lib/bootstrap.min.css), [static/css/lib/tempus-dominus.min.css](../../static/css/lib/tempus-dominus.min.css), [static/js/cheesy-common.js](../../static/js/cheesy-common.js), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/lib/bootstrap.bundle.min.js](../../static/js/lib/bootstrap.bundle.min.js), [static/js/lib/handlebars-1.3.0.js](../../static/js/lib/handlebars-1.3.0.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/lib/moment.min.js](../../static/js/lib/moment.min.js), [static/js/lib/popper.min.js](../../static/js/lib/popper.min.js), [static/js/lib/tempus-dominus.min.js](../../static/js/lib/tempus-dominus.min.js) |
| [templates/bracket.svg](../../templates/bracket.svg) |  |
| [templates/bracket_display.html](../../templates/bracket_display.html) | [static/css/bracket_display.css](../../static/css/bracket_display.css), [static/css/cheesy-arena.css](../../static/css/cheesy-arena.css), [static/css/lib/bootstrap.min.css](../../static/css/lib/bootstrap.min.css), [static/js/bracket_display.js](../../static/js/bracket_display.js), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js) |
| [templates/bracket_report.html](../../templates/bracket_report.html) |  |
| [templates/edit_match_result.html](../../templates/edit_match_result.html) | [static/js/match_review.js](../../static/js/match_review.js) |
| [templates/edit_team.html](../../templates/edit_team.html) |  |
| [templates/field_monitor_display.html](../../templates/field_monitor_display.html) | [static/css/cheesy-arena.css](../../static/css/cheesy-arena.css), [static/css/field_monitor_display.css](../../static/css/field_monitor_display.css), [static/css/lib/bootstrap-icons.min.css](../../static/css/lib/bootstrap-icons.min.css), [static/css/lib/bootstrap.min.css](../../static/css/lib/bootstrap.min.css), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/field_monitor_display.js](../../static/js/field_monitor_display.js), [static/js/lib/bootstrap.bundle.min.js](../../static/js/lib/bootstrap.bundle.min.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.transit.min.js](../../static/js/lib/jquery.transit.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/match_timing.js](../../static/js/match_timing.js) |
| [templates/fms_field_monitor_display.html](../../templates/fms_field_monitor_display.html) | [static/css/cheesy-arena.css](../../static/css/cheesy-arena.css), [static/css/fms_field_monitor_display.css](../../static/css/fms_field_monitor_display.css), [static/css/lib/bootstrap-icons.min.css](../../static/css/lib/bootstrap-icons.min.css), [static/css/lib/bootstrap.min.css](../../static/css/lib/bootstrap.min.css), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/fms_field_monitor_display.js](../../static/js/fms_field_monitor_display.js), [static/js/lib/bootstrap.bundle.min.js](../../static/js/lib/bootstrap.bundle.min.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.transit.min.js](../../static/js/lib/jquery.transit.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/match_timing.js](../../static/js/match_timing.js) |
| [templates/fta.csv](../../templates/fta.csv) |  |
| [templates/index.html](../../templates/index.html) |  |
| [templates/login.html](../../templates/login.html) |  |
| [templates/logo_display.html](../../templates/logo_display.html) | [static/css/cheesy-arena.css](../../static/css/cheesy-arena.css), [static/css/lib/bootstrap.min.css](../../static/css/lib/bootstrap.min.css), [static/css/logo_display.css](../../static/css/logo_display.css), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/logo_display.js](../../static/js/logo_display.js), [static/js/match_timing.js](../../static/js/match_timing.js) |
| [templates/match_logs.html](../../templates/match_logs.html) |  |
| [templates/match_play.html](../../templates/match_play.html) | [static/js/match_play.js](../../static/js/match_play.js), [static/js/match_timing.js](../../static/js/match_timing.js) |
| [templates/match_play_match_load.html](../../templates/match_play_match_load.html) |  |
| [templates/match_review.html](../../templates/match_review.html) |  |
| [templates/placeholder_display.html](../../templates/placeholder_display.html) | [static/css/cheesy-arena.css](../../static/css/cheesy-arena.css), [static/css/placeholder_display.css](../../static/css/placeholder_display.css), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/placeholder_display.js](../../static/js/placeholder_display.js) |
| [templates/queueing_display.html](../../templates/queueing_display.html) | [static/css/cheesy-arena.css](../../static/css/cheesy-arena.css), [static/css/lib/bootstrap.min.css](../../static/css/lib/bootstrap.min.css), [static/css/queueing_display.css](../../static/css/queueing_display.css), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/lib/bootstrap.bundle.min.js](../../static/js/lib/bootstrap.bundle.min.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.transit.min.js](../../static/js/lib/jquery.transit.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/match_timing.js](../../static/js/match_timing.js), [static/js/queueing_display.js](../../static/js/queueing_display.js) |
| [templates/queueing_display_match_load.html](../../templates/queueing_display_match_load.html) |  |
| [templates/rankings.csv](../../templates/rankings.csv) |  |
| [templates/rankings_display.html](../../templates/rankings_display.html) | [static/css/cheesy-arena.css](../../static/css/cheesy-arena.css), [static/css/lib/bootstrap.min.css](../../static/css/lib/bootstrap.min.css), [static/css/rankings_display.css](../../static/css/rankings_display.css), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/lib/handlebars-1.3.0.js](../../static/js/lib/handlebars-1.3.0.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.transit.min.js](../../static/js/lib/jquery.transit.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/rankings_display.js](../../static/js/rankings_display.js) |
| [templates/referee_panel.html](../../templates/referee_panel.html) | [static/css/referee_panel.css](../../static/css/referee_panel.css), [static/js/match_timing.js](../../static/js/match_timing.js), [static/js/referee_panel.js](../../static/js/referee_panel.js) |
| [templates/referee_panel_foul_list.html](../../templates/referee_panel_foul_list.html) |  |
| [templates/schedule.csv](../../templates/schedule.csv) |  |
| [templates/scoring_panel.html](../../templates/scoring_panel.html) | [static/css/scoring_panel.css](../../static/css/scoring_panel.css), [static/js/match_timing.js](../../static/js/match_timing.js), [static/js/scoring_panel.js](../../static/js/scoring_panel.js) |
| [templates/setup_awards.html](../../templates/setup_awards.html) |  |
| [templates/setup_breaks.html](../../templates/setup_breaks.html) |  |
| [templates/setup_displays.html](../../templates/setup_displays.html) | [static/js/setup_displays.js](../../static/js/setup_displays.js) |
| [templates/setup_field_testing.html](../../templates/setup_field_testing.html) | [static/js/setup_field_testing.js](../../static/js/setup_field_testing.js) |
| [templates/setup_judging.html](../../templates/setup_judging.html) |  |
| [templates/setup_lower_thirds.html](../../templates/setup_lower_thirds.html) | [static/js/lower_thirds.js](../../static/js/lower_thirds.js) |
| [templates/setup_schedule.html](../../templates/setup_schedule.html) | [static/js/setup_schedule.js](../../static/js/setup_schedule.js) |
| [templates/setup_settings.html](../../templates/setup_settings.html) |  |
| [templates/setup_sponsor_slides.html](../../templates/setup_sponsor_slides.html) |  |
| [templates/setup_teams.html](../../templates/setup_teams.html) |  |
| [templates/teams.csv](../../templates/teams.csv) |  |
| [templates/twitch_display.html](../../templates/twitch_display.html) | [static/css/twitch_display.css](../../static/css/twitch_display.css), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/twitch_display.js](../../static/js/twitch_display.js) |
| [templates/unpicked_display.html](../../templates/unpicked_display.html) | [static/css/cheesy-arena.css](../../static/css/cheesy-arena.css), [static/css/lib/bootstrap.min.css](../../static/css/lib/bootstrap.min.css), [static/css/unpicked_display.css](../../static/css/unpicked_display.css), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/lib/handlebars-1.3.0.js](../../static/js/lib/handlebars-1.3.0.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/unpicked_display.js](../../static/js/unpicked_display.js) |
| [templates/view_match_log.html](../../templates/view_match_log.html) |  |
| [templates/wall_display.html](../../templates/wall_display.html) | [static/css/cheesy-arena.css](../../static/css/cheesy-arena.css), [static/css/lib/bootstrap-icons.min.css](../../static/css/lib/bootstrap-icons.min.css), [static/css/lib/bootstrap.min.css](../../static/css/lib/bootstrap.min.css), [static/css/wall_display.css](../../static/css/wall_display.css), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/lib/bootstrap.bundle.min.js](../../static/js/lib/bootstrap.bundle.min.js), [static/js/lib/handlebars-1.3.0.js](../../static/js/lib/handlebars-1.3.0.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.transit.min.js](../../static/js/lib/jquery.transit.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/match_timing.js](../../static/js/match_timing.js), [static/js/wall_display.js](../../static/js/wall_display.js) |
| [templates/webpage_display.html](../../templates/webpage_display.html) | [static/css/webpage_display.css](../../static/css/webpage_display.css), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/webpage_display.js](../../static/js/webpage_display.js) |
| [templates/youtube_display.html](../../templates/youtube_display.html) | [static/css/youtube_display.css](../../static/css/youtube_display.css), [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js), [static/js/lib/jquery.json-2.4.min.js](../../static/js/lib/jquery.json-2.4.min.js), [static/js/lib/jquery.min.js](../../static/js/lib/jquery.min.js), [static/js/lib/jquery.websocket-0.0.1.js](../../static/js/lib/jquery.websocket-0.0.1.js), [static/js/youtube_display.js](../../static/js/youtube_display.js) |

## Uygulamaya ait tarayıcı kodu

| Dosya | Satır | WebSocket uçları |
| --- | --- | --- |
| [static/js/alliance_selection.js](../../static/js/alliance_selection.js) | 64 | `/alliance_selection/websocket` |
| [static/js/alliance_station_display.js](../../static/js/alliance_station_display.js) | 157 | `/displays/alliance_station/websocket` |
| [static/js/announcer_display.js](../../static/js/announcer_display.js) | 106 | `/displays/announcer/websocket` |
| [static/js/audience_display.js](../../static/js/audience_display.js) | 850 | `/displays/audience/websocket` |
| [static/js/bracket_display.js](../../static/js/bracket_display.js) | 22 | `/displays/bracket/websocket` |
| [static/js/cheesy-common.js](../../static/js/cheesy-common.js) | 34 |  |
| [static/js/cheesy-websocket.js](../../static/js/cheesy-websocket.js) | 76 |  |
| [static/js/display_shared.js](../../static/js/display_shared.js) | 275 |  |
| [static/js/field_monitor_display.js](../../static/js/field_monitor_display.js) | 266 | `/displays/field_monitor/websocket` |
| [static/js/fms_field_monitor_display.js](../../static/js/fms_field_monitor_display.js) | 445 | `/displays/fms_field_monitor/websocket` |
| [static/js/logo_display.js](../../static/js/logo_display.js) | 18 | `/displays/logo/websocket` |
| [static/js/lower_thirds.js](../../static/js/lower_thirds.js) | 66 | `/setup/lower_thirds/websocket` |
| [static/js/match_play.js](../../static/js/match_play.js) | 463 | `/match_play/websocket` |
| [static/js/match_review.js](../../static/js/match_review.js) | 258 |  |
| [static/js/match_timing.js](../../static/js/match_timing.js) | 86 |  |
| [static/js/placeholder_display.js](../../static/js/placeholder_display.js) | 19 | `/display/websocket` |
| [static/js/queueing_display.js](../../static/js/queueing_display.js) | 49 | `/displays/queueing/websocket` |
| [static/js/rankings_display.js](../../static/js/rankings_display.js) | 100 | `/displays/rankings/websocket` |
| [static/js/referee_panel.js](../../static/js/referee_panel.js) | 291 | `/panels/referee/websocket` |
| [static/js/scoring_panel.js](../../static/js/scoring_panel.js) | 186 | `/websocket` |
| [static/js/setup_displays.js](../../static/js/setup_displays.js) | 284 | `/setup/displays/websocket` |
| [static/js/setup_field_testing.js](../../static/js/setup_field_testing.js) | 176 | `/setup/field_testing/websocket` |
| [static/js/setup_schedule.js](../../static/js/setup_schedule.js) | 113 |  |
| [static/js/stream_displays_test.js](../../static/js/stream_displays_test.js) | 91 | `/displays/twitch/websocket`, `/displays/youtube/websocket` |
| [static/js/twitch_display.js](../../static/js/twitch_display.js) | 30 | `/displays/twitch/websocket` |
| [static/js/unpicked_display.js](../../static/js/unpicked_display.js) | 47 | `/displays/unpicked/websocket` |
| [static/js/wall_display.js](../../static/js/wall_display.js) | 357 | `/displays/wall/websocket` |
| [static/js/webpage_display.js](../../static/js/webpage_display.js) | 15 | `/displays/webpage/websocket` |
| [static/js/youtube_display.js](../../static/js/youtube_display.js) | 43 | `/displays/youtube/websocket` |

## Çalışma zamanı varlıkları

| Tür | Git içindeki dosya sayısı |
| --- | --- |
| schedules/*.csv | 1330 |
| templates/* | 50 |
| static/js (lib hariç) | 29 |
| static/css (lib ve fonts hariç) | 18 |
| static/audio/* | 10 |
