# WebSocket sözleşmeleri

Üretilen envanter; veri alanlarının anlamı ana sistem haritasında açıklanır. Çıkarım, handler içindeki HandleNotifiers çağrıları ve üst seviye komut dalları üzerinden yapılır. `display.Notifier` dinamik ekran URL yapılandırmasıdır. `reload`, `playSound` gibi üreticisi nil olaylar bağlantı başında gönderilmez.

## Mesaj zarfı

```json
{"type": "startMatch", "data": {"MuteMatchSounds": false}}
```

Alan isimleri veri gövdelerindeki Go struct adlarına göre çoğunlukla PascalCase biçimindedir. `ping` ve `error` altyapı mesajlarıdır. Her bağlantı kendi notifier kanallarına abone olur; kanallar arası toplam sıra garantisi yoktur.

## Arena olayları

| Arena alanı | JSON type | Veri üreticisi |
| --- | --- | --- |
| `AllianceSelectionNotifier` | `allianceSelection` | `arena.generateAllianceSelectionMessage` |
| `AllianceStationDisplayModeNotifier` | `allianceStationDisplayMode` | `arena.generateAllianceStationDisplayModeMessage` |
| `ArenaStatusNotifier` | `arenaStatus` | `arena.generateArenaStatusMessage` |
| `AudienceDisplayModeNotifier` | `audienceDisplayMode` | `arena.generateAudienceDisplayModeMessage` |
| `DisplayConfigurationNotifier` | `displayConfiguration` | `arena.generateDisplayConfigurationMessage` |
| `EventStatusNotifier` | `eventStatus` | `arena.generateEventStatusMessage` |
| `LowerThirdNotifier` | `lowerThird` | `arena.generateLowerThirdMessage` |
| `MatchLoadNotifier` | `matchLoad` | `arena.GenerateMatchLoadMessage` |
| `MatchTimeNotifier` | `matchTime` | `arena.generateMatchTimeMessage` |
| `MatchTimingNotifier` | `matchTiming` | `arena.generateMatchTimingMessage` |
| `PlaySoundNotifier` | `playSound` | `nil` |
| `RealtimeScoreNotifier` | `realtimeScore` | `arena.generateRealtimeScoreMessage` |
| `ReloadDisplaysNotifier` | `reload` | `nil` |
| `ScorePostedNotifier` | `scorePosted` | `arena.GenerateScorePostedMessage` |
| `ScoringStatusNotifier` | `scoringStatus` | `arena.generateScoringStatusMessage` |

## `/alliance_selection/websocket`

[allianceSelectionWebsocketHandler](../../web/alliance_selection.go#L274)

- Admin kontrolü: var.
- Abonelikler: `AllianceSelectionNotifier`, `AudienceDisplayModeNotifier`.
- İstemci komutları: `setTimer`, `startTimer`, `stopTimer`, `restartTimer`, `hideTimer`, `setAudienceDisplay`.

## `/api/arena/websocket`

[arenaWebsocketApiHandler](../../web/api.go#L208)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `MatchTimingNotifier`, `MatchLoadNotifier`, `MatchTimeNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/display/websocket`

[placeholderDisplayWebsocketHandler](../../web/placeholder_display.go#L36)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `display.Notifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/displays/alliance_station/websocket`

[allianceStationDisplayWebsocketHandler](../../web/alliance_station_display.go#L37)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `display.Notifier`, `MatchTimingNotifier`, `AllianceStationDisplayModeNotifier`, `ArenaStatusNotifier`, `MatchLoadNotifier`, `MatchTimeNotifier`, `RealtimeScoreNotifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/displays/announcer/websocket`

[announcerDisplayWebsocketHandler](../../web/announcer_display.go#L67)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `display.Notifier`, `MatchTimingNotifier`, `AudienceDisplayModeNotifier`, `EventStatusNotifier`, `MatchLoadNotifier`, `MatchTimeNotifier`, `RealtimeScoreNotifier`, `ScorePostedNotifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/displays/audience/websocket`

[audienceDisplayWebsocketHandler](../../web/audience_display.go#L48)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `display.Notifier`, `MatchTimingNotifier`, `AudienceDisplayModeNotifier`, `MatchLoadNotifier`, `MatchTimeNotifier`, `RealtimeScoreNotifier`, `PlaySoundNotifier`, `ScorePostedNotifier`, `AllianceSelectionNotifier`, `LowerThirdNotifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/displays/bracket/websocket`

[bracketDisplayWebsocketHandler](../../web/bracket_display.go#L36)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `display.Notifier`, `MatchLoadNotifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/displays/field_monitor/websocket`

[fieldMonitorDisplayWebsocketHandler](../../web/field_monitor_display.go#L68)

- Admin kontrolü: yalnızca fta=true ise.
- Abonelikler: `MatchTimingNotifier`, `display.Notifier`, `ArenaStatusNotifier`, `EventStatusNotifier`, `RealtimeScoreNotifier`, `MatchTimeNotifier`, `MatchLoadNotifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: `updateTeamNotes`.

## `/displays/fms_field_monitor/websocket`

[fieldMonitorDisplayWebsocketHandler](../../web/field_monitor_display.go#L68)

- Admin kontrolü: yalnızca fta=true ise.
- Abonelikler: `MatchTimingNotifier`, `display.Notifier`, `ArenaStatusNotifier`, `EventStatusNotifier`, `RealtimeScoreNotifier`, `MatchTimeNotifier`, `MatchLoadNotifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: `updateTeamNotes`.

## `/displays/logo/websocket`

[logoDisplayWebsocketHandler](../../web/logo_display.go#L36)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `display.Notifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/displays/queueing/websocket`

[queueingDisplayWebsocketHandler](../../web/queueing_display.go#L110)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `display.Notifier`, `MatchTimingNotifier`, `MatchLoadNotifier`, `MatchTimeNotifier`, `EventStatusNotifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/displays/rankings/websocket`

[rankingsDisplayWebsocketHandler](../../web/rankings_display.go#L36)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `display.Notifier`, `EventStatusNotifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/displays/twitch/websocket`

[twitchDisplayWebsocketHandler](../../web/twitch_display.go#L36)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `display.Notifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/displays/unpicked/websocket`

[unpickedDisplayWebsocketHandler](../../web/unpicked_display.go#L38)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `display.Notifier`, `AllianceSelectionNotifier`, `AudienceDisplayModeNotifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/displays/wall/websocket`

[wallDisplayWebsocketHandler](../../web/wall_display.go#L49)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `display.Notifier`, `MatchTimingNotifier`, `AudienceDisplayModeNotifier`, `MatchLoadNotifier`, `MatchTimeNotifier`, `RealtimeScoreNotifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/displays/webpage/websocket`

[webpageDisplayWebsocketHandler](../../web/webpage_display.go#L36)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `display.Notifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/displays/youtube/websocket`

[youtubeDisplayWebsocketHandler](../../web/youtube_display.go#L36)

- Admin kontrolü: doğrudan yok.
- Abonelikler: `display.Notifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: handler içinde komut okuma döngüsü yok.

## `/match_play/websocket`

[matchPlayWebsocketHandler](../../web/match_play.go#L116)

- Admin kontrolü: var.
- Abonelikler: `MatchTimingNotifier`, `AllianceStationDisplayModeNotifier`, `ArenaStatusNotifier`, `AudienceDisplayModeNotifier`, `EventStatusNotifier`, `MatchLoadNotifier`, `MatchTimeNotifier`, `RealtimeScoreNotifier`, `ScorePostedNotifier`, `ScoringStatusNotifier`.
- İstemci komutları: `loadMatch`, `showResult`, `substituteTeams`, `toggleBypass`, `startMatch`, `abortMatch`, `signalVolunteers`, `signalReset`, `commitAndPost`, `discardResults`, `setAudienceDisplay`, `setAllianceStationDisplay`, `startTimeout`, `setTimeoutDisplay`, `setTestMatchName`.

## `/panels/scoring/{position}/websocket`

[scoringPanelWebsocketHandler](../../web/scoring_panel.go#L67)

- Admin kontrolü: var.
- Abonelikler: `MatchLoadNotifier`, `MatchTimeNotifier`, `RealtimeScoreNotifier`, `ReloadDisplaysNotifier`.
- İstemci komutları: `commitMatch`, `autoTower`, `endgame`, `addFoul`.
- Doğrudan yanıtlar: `resetLocalState`.

## `/panels/referee/websocket`

[refereePanelWebsocketHandler](../../web/referee_panel.go#L71)

- Admin kontrolü: var.
- Abonelikler: `MatchLoadNotifier`, `MatchTimeNotifier`, `RealtimeScoreNotifier`, `ScoringStatusNotifier`, `ReloadDisplaysNotifier`, `ArenaStatusNotifier`.
- İstemci komutları: `addFoul`, `toggleFoulType`, `updateFoulTeam`, `updateFoulRule`, `deleteFoul`, `card`, `toggleBypass`, `signalVolunteers`, `signalReset`, `commitAndPost`.

## `/setup/displays/websocket`

[displaysWebsocketHandler](../../web/setup_displays.go#L42)

- Admin kontrolü: var.
- Abonelikler: `DisplayConfigurationNotifier`.
- İstemci komutları: `configureDisplay`, `reloadDisplay`, `reloadAllDisplays`.
- Doğrudan yanıtlar: `displayConfigurationSaved`.

## `/setup/field_testing/websocket`

[fieldTestingWebsocketHandler](../../web/setup_field_testing.go#L69)

- Admin kontrolü: var.
- Abonelikler: `Plc.IoChangeNotifier()`, `ArenaStatusNotifier`.
- İstemci komutları: `playSound`, `setPlcCoilOverride`, `setLedMode`.
- Doğrudan yanıtlar: `ledStatus`.

## `/setup/lower_thirds/websocket`

[lowerThirdsWebsocketHandler](../../web/setup_lower_thirds.go#L48)

- Admin kontrolü: var.
- Abonelikler: `AudienceDisplayModeNotifier`.
- İstemci komutları: `saveLowerThird`, `deleteLowerThird`, `showLowerThird`, `hideLowerThird`, `reorderLowerThird`, `setAudienceDisplay`.
