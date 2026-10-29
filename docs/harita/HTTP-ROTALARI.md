# HTTP rota envanteri

Kaynak: `web/web.go` içindeki açık HandleFunc kayıtları. Admin sütunu ilgili handler gövdesindeki doğrudan `userIsAdmin` çağrısını gösterir; bağımsız bir güvenlik denetimi değildir. Admin parolası boşsa bu kontrol de erişimi açık bırakır.

Toplam **115** açık rota; **23** WebSocket rotası. Ek olarak `/static/` dosya sunucusu ve `/` yönlendiricisi ServeWebInterface içinde kaydedilir. GET kalıpları Go ServeMux davranışı gereği HEAD isteklerini de karşılayabilir.

| Yöntem | Yol | Handler / kaynak | Erişim kontrolü |
| --- | --- | --- | --- |
| GET | `/` | [indexHandler](../../web/web.go#L101) | Doğrudan admin kontrolü yok |
| GET | `/alliance_selection` | [allianceSelectionGetHandler](../../web/alliance_selection.go#L33) | Admin kontrolü |
| POST | `/alliance_selection` | [allianceSelectionPostHandler](../../web/alliance_selection.go#L42) | Admin kontrolü |
| GET | `/alliance_selection/websocket` | [allianceSelectionWebsocketHandler](../../web/alliance_selection.go#L274) | Admin kontrolü |
| POST | `/alliance_selection/finalize` | [allianceSelectionFinalizeHandler](../../web/alliance_selection.go#L178) | Admin kontrolü |
| POST | `/alliance_selection/reset` | [allianceSelectionResetHandler](../../web/alliance_selection.go#L148) | Admin kontrolü |
| POST | `/alliance_selection/start` | [allianceSelectionStartHandler](../../web/alliance_selection.go#L103) | Admin kontrolü |
| GET | `/api/alliances` | [alliancesApiHandler](../../web/api.go#L185) | Doğrudan admin kontrolü yok |
| GET | `/api/arena/websocket` | [arenaWebsocketApiHandler](../../web/api.go#L208) | Doğrudan admin kontrolü yok |
| GET | `/api/bracket/svg` | [bracketSvgApiHandler](../../web/api.go#L236) | Doğrudan admin kontrolü yok |
| GET | `/api/matches/{type}` | [matchesApiHandler](../../web/api.go#L51) | Doğrudan admin kontrolü yok |
| GET | `/api/rankings` | [rankingsApiHandler](../../web/api.go#L124) | Doğrudan admin kontrolü yok |
| GET | `/api/sponsor_slides` | [sponsorSlidesApiHandler](../../web/api.go#L97) | Doğrudan admin kontrolü yok |
| GET | `/api/teams/{teamId}/avatar` | [teamAvatarsApiHandler](../../web/api.go#L221) | Doğrudan admin kontrolü yok |
| GET | `/display` | [placeholderDisplayHandler](../../web/placeholder_display.go#L15) | Doğrudan admin kontrolü yok |
| GET | `/display/websocket` | [placeholderDisplayWebsocketHandler](../../web/placeholder_display.go#L36) | Doğrudan admin kontrolü yok |
| GET | `/displays/alliance_station` | [allianceStationDisplayHandler](../../web/alliance_station_display.go#L15) | Doğrudan admin kontrolü yok |
| GET | `/displays/alliance_station/websocket` | [allianceStationDisplayWebsocketHandler](../../web/alliance_station_display.go#L37) | Doğrudan admin kontrolü yok |
| GET | `/displays/announcer` | [announcerDisplayHandler](../../web/announcer_display.go#L15) | Doğrudan admin kontrolü yok |
| GET | `/displays/announcer/match_load` | [announcerDisplayMatchLoadHandler](../../web/announcer_display.go#L37) | Doğrudan admin kontrolü yok |
| GET | `/displays/announcer/score_posted` | [announcerDisplayScorePostedHandler](../../web/announcer_display.go#L52) | Doğrudan admin kontrolü yok |
| GET | `/displays/announcer/websocket` | [announcerDisplayWebsocketHandler](../../web/announcer_display.go#L67) | Doğrudan admin kontrolü yok |
| GET | `/displays/audience` | [audienceDisplayHandler](../../web/audience_display.go#L16) | Doğrudan admin kontrolü yok |
| GET | `/displays/audience/websocket` | [audienceDisplayWebsocketHandler](../../web/audience_display.go#L48) | Doğrudan admin kontrolü yok |
| GET | `/displays/bracket` | [bracketDisplayHandler](../../web/bracket_display.go#L15) | Doğrudan admin kontrolü yok |
| GET | `/displays/bracket/websocket` | [bracketDisplayWebsocketHandler](../../web/bracket_display.go#L36) | Doğrudan admin kontrolü yok |
| GET | `/displays/field_monitor` | [fieldMonitorDisplayHandler](../../web/field_monitor_display.go#L18) | Yalnızca fta=true ise admin kontrolü |
| GET | `/displays/field_monitor/websocket` | [fieldMonitorDisplayWebsocketHandler](../../web/field_monitor_display.go#L68) | Yalnızca fta=true ise admin kontrolü |
| GET | `/displays/fms_field_monitor` | [fmsFieldMonitorDisplayHandler](../../web/field_monitor_display.go#L43) | Yalnızca fta=true ise admin kontrolü |
| GET | `/displays/fms_field_monitor/websocket` | [fieldMonitorDisplayWebsocketHandler](../../web/field_monitor_display.go#L68) | Yalnızca fta=true ise admin kontrolü |
| GET | `/displays/logo` | [logoDisplayHandler](../../web/logo_display.go#L15) | Doğrudan admin kontrolü yok |
| GET | `/displays/logo/websocket` | [logoDisplayWebsocketHandler](../../web/logo_display.go#L36) | Doğrudan admin kontrolü yok |
| GET | `/displays/queueing` | [queueingDisplayHandler](../../web/queueing_display.go#L22) | Doğrudan admin kontrolü yok |
| GET | `/displays/queueing/match_load` | [queueingDisplayMatchLoadHandler](../../web/queueing_display.go#L46) | Doğrudan admin kontrolü yok |
| GET | `/displays/queueing/websocket` | [queueingDisplayWebsocketHandler](../../web/queueing_display.go#L110) | Doğrudan admin kontrolü yok |
| GET | `/displays/rankings` | [rankingsDisplayHandler](../../web/rankings_display.go#L15) | Doğrudan admin kontrolü yok |
| GET | `/displays/rankings/websocket` | [rankingsDisplayWebsocketHandler](../../web/rankings_display.go#L36) | Doğrudan admin kontrolü yok |
| GET | `/displays/twitch` | [twitchDisplayHandler](../../web/twitch_display.go#L15) | Doğrudan admin kontrolü yok |
| GET | `/displays/twitch/websocket` | [twitchDisplayWebsocketHandler](../../web/twitch_display.go#L36) | Doğrudan admin kontrolü yok |
| GET | `/displays/unpicked` | [unpickedDisplayHandler](../../web/unpicked_display.go#L13) | Doğrudan admin kontrolü yok |
| GET | `/displays/unpicked/websocket` | [unpickedDisplayWebsocketHandler](../../web/unpicked_display.go#L38) | Doğrudan admin kontrolü yok |
| GET | `/displays/wall` | [wallDisplayHandler](../../web/wall_display.go#L16) | Doğrudan admin kontrolü yok |
| GET | `/displays/wall/websocket` | [wallDisplayWebsocketHandler](../../web/wall_display.go#L49) | Doğrudan admin kontrolü yok |
| GET | `/displays/webpage` | [webpageDisplayHandler](../../web/webpage_display.go#L15) | Doğrudan admin kontrolü yok |
| GET | `/displays/webpage/websocket` | [webpageDisplayWebsocketHandler](../../web/webpage_display.go#L36) | Doğrudan admin kontrolü yok |
| GET | `/displays/youtube` | [youtubeDisplayHandler](../../web/youtube_display.go#L15) | Doğrudan admin kontrolü yok |
| GET | `/displays/youtube/websocket` | [youtubeDisplayWebsocketHandler](../../web/youtube_display.go#L36) | Doğrudan admin kontrolü yok |
| GET | `/login` | [loginHandler](../../web/login.go#L19) | Doğrudan admin kontrolü yok |
| POST | `/login` | [loginPostHandler](../../web/login.go#L24) | Doğrudan admin kontrolü yok |
| GET | `/match_play` | [matchPlayHandler](../../web/match_play.go#L36) | Admin kontrolü |
| GET | `/match_play/match_load` | [matchPlayMatchLoadHandler](../../web/match_play.go#L65) | Admin kontrolü |
| GET | `/match_play/websocket` | [matchPlayWebsocketHandler](../../web/match_play.go#L116) | Admin kontrolü |
| GET | `/match_logs` | [matchLogsHandler](../../web/match_logs.go#L70) | Doğrudan admin kontrolü yok |
| GET | `/match_logs/{matchId}/{stationId}/log` | [matchLogsViewGetHandler](../../web/match_logs.go#L114) | Doğrudan admin kontrolü yok |
| GET | `/match_review` | [matchReviewHandler](../../web/match_review.go#L45) | Doğrudan admin kontrolü yok |
| GET | `/match_review/{matchId}/edit` | [matchReviewEditGetHandler](../../web/match_review.go#L89) | Admin kontrolü |
| POST | `/match_review/{matchId}/edit` | [matchReviewEditPostHandler](../../web/match_review.go#L183) | Admin kontrolü |
| POST | `/match_review/{matchId}/summary` | [matchReviewSummaryPostHandler](../../web/match_review.go#L142) | Admin kontrolü |
| GET | `/panels/scoring/{position}` | [scoringPanelHandler](../../web/scoring_panel.go#L37) | Admin kontrolü |
| GET | `/panels/scoring/{position}/websocket` | [scoringPanelWebsocketHandler](../../web/scoring_panel.go#L67) | Admin kontrolü |
| GET | `/panels/referee` | [refereePanelHandler](../../web/referee_panel.go#L23) | Admin kontrolü |
| GET | `/panels/referee/foul_list` | [refereePanelFoulListHandler](../../web/referee_panel.go#L45) | Doğrudan admin kontrolü yok |
| GET | `/panels/referee/websocket` | [refereePanelWebsocketHandler](../../web/referee_panel.go#L71) | Admin kontrolü |
| GET | `/reports/csv/backups` | [backupTeamsCsvReportHandler](../../web/reports.go#L173) | Doğrudan admin kontrolü yok |
| GET | `/reports/csv/fta` | [ftaCsvReportHandler](../../web/reports.go#L949) | Doğrudan admin kontrolü yok |
| GET | `/reports/csv/rankings` | [rankingsCsvReportHandler](../../web/reports.go#L25) | Doğrudan admin kontrolü yok |
| GET | `/reports/csv/schedule/{type}` | [scheduleCsvReportHandler](../../web/reports.go#L374) | Doğrudan admin kontrolü yok |
| GET | `/reports/csv/teams` | [teamsCsvReportHandler](../../web/reports.go#L554) | Doğrudan admin kontrolü yok |
| GET | `/reports/csv/wpa_keys` | [wpaKeysCsvReportHandler](../../web/reports.go#L661) | Admin kontrolü |
| GET | `/reports/pdf/alliances` | [alliancesPdfReportHandler](../../web/reports.go#L684) | Doğrudan admin kontrolü yok |
| GET | `/reports/pdf/backups` | [backupsPdfReportHandler](../../web/reports.go#L223) | Doğrudan admin kontrolü yok |
| GET | `/reports/pdf/bracket` | [bracketPdfReportHandler](../../web/reports.go#L807) | Doğrudan admin kontrolü yok |
| GET | `/reports/pdf/coupons` | [couponsPdfReportHandler](../../web/reports.go#L290) | Doğrudan admin kontrolü yok |
| GET | `/reports/pdf/cycle/{type}` | [cyclePdfReportHandler](../../web/reports.go#L845) | Doğrudan admin kontrolü yok |
| GET | `/reports/pdf/judging_schedule` | [judgingSchedulePdfReportHandler](../../web/reports.go#L979) | Doğrudan admin kontrolü yok |
| GET | `/reports/pdf/rankings` | [rankingsPdfReportHandler](../../web/reports.go#L55) | Doğrudan admin kontrolü yok |
| GET | `/reports/pdf/schedule/{type}` | [schedulePdfReportHandler](../../web/reports.go#L410) | Doğrudan admin kontrolü yok |
| GET | `/reports/pdf/teams` | [teamsPdfReportHandler](../../web/reports.go#L584) | Doğrudan admin kontrolü yok |
| GET | `/setup/awards` | [awardsGetHandler](../../web/setup_awards.go#L17) | Admin kontrolü |
| POST | `/setup/awards` | [awardsPostHandler](../../web/setup_awards.go#L54) | Admin kontrolü |
| GET | `/setup/breaks` | [breaksGetHandler](../../web/setup_breaks.go#L15) | Admin kontrolü |
| POST | `/setup/breaks` | [breaksPostHandler](../../web/setup_breaks.go#L42) | Admin kontrolü |
| POST | `/setup/db/clear/{type}` | [clearDbHandler](../../web/setup_settings.go#L330) | Admin kontrolü |
| POST | `/setup/db/restore` | [restoreDbHandler](../../web/setup_settings.go#L244) | Admin kontrolü |
| GET | `/setup/db/save` | [saveDbHandler](../../web/setup_settings.go#L227) | Admin kontrolü |
| GET | `/setup/displays` | [displaysGetHandler](../../web/setup_displays.go#L20) | Admin kontrolü |
| GET | `/setup/displays/websocket` | [displaysWebsocketHandler](../../web/setup_displays.go#L42) | Admin kontrolü |
| GET | `/setup/field_testing` | [fieldTestingGetHandler](../../web/setup_field_testing.go#L28) | Admin kontrolü |
| GET | `/setup/field_testing/websocket` | [fieldTestingWebsocketHandler](../../web/setup_field_testing.go#L69) | Admin kontrolü |
| GET | `/setup/judging` | [judgingGetHandler](../../web/setup_judging.go#L25) | Admin kontrolü |
| POST | `/setup/judging/clear` | [judgingClearPostHandler](../../web/setup_judging.go#L92) | Admin kontrolü |
| POST | `/setup/judging/generate` | [judgingGeneratePostHandler](../../web/setup_judging.go#L34) | Admin kontrolü |
| GET | `/setup/lower_thirds` | [lowerThirdsGetHandler](../../web/setup_lower_thirds.go#L19) | Admin kontrolü |
| GET | `/setup/lower_thirds/websocket` | [lowerThirdsWebsocketHandler](../../web/setup_lower_thirds.go#L48) | Admin kontrolü |
| GET | `/setup/schedule` | [scheduleGetHandler](../../web/setup_schedule.go#L22) | Admin kontrolü |
| POST | `/setup/schedule/generate` | [scheduleGeneratePostHandler](../../web/setup_schedule.go#L38) | Admin kontrolü |
| POST | `/setup/schedule/save` | [scheduleSavePostHandler](../../web/setup_schedule.go#L122) | Admin kontrolü |
| GET | `/setup/settings` | [settingsGetHandler](../../web/setup_settings.go#L23) | Admin kontrolü |
| POST | `/setup/settings` | [settingsPostHandler](../../web/setup_settings.go#L32) | Admin kontrolü |
| GET | `/setup/settings/publish_alliances` | [settingsPublishAlliancesHandler](../../web/setup_settings.go#L381) | Admin kontrolü |
| GET | `/setup/settings/publish_awards` | [settingsPublishAwardsHandler](../../web/setup_settings.go#L403) | Admin kontrolü |
| GET | `/setup/settings/publish_matches` | [settingsPublishMatchesHandler](../../web/setup_settings.go#L425) | Admin kontrolü |
| GET | `/setup/settings/publish_rankings` | [settingsPublishRankingsHandler](../../web/setup_settings.go#L454) | Admin kontrolü |
| GET | `/setup/settings/publish_teams` | [settingsPublishTeamsHandler](../../web/setup_settings.go#L476) | Admin kontrolü |
| GET | `/setup/sponsor_slides` | [sponsorSlidesGetHandler](../../web/setup_sponsor_slides.go#L15) | Admin kontrolü |
| POST | `/setup/sponsor_slides` | [sponsorSlidesPostHandler](../../web/setup_sponsor_slides.go#L46) | Admin kontrolü |
| GET | `/setup/teams` | [teamsGetHandler](../../web/setup_teams.go#L26) | Admin kontrolü |
| POST | `/setup/teams` | [teamsPostHandler](../../web/setup_teams.go#L35) | Admin kontrolü |
| POST | `/setup/teams/{id}/delete` | [teamDeletePostHandler](../../web/setup_teams.go#L202) | Admin kontrolü |
| GET | `/setup/teams/{id}/edit` | [teamEditGetHandler](../../web/setup_teams.go#L127) | Admin kontrolü |
| POST | `/setup/teams/{id}/edit` | [teamEditPostHandler](../../web/setup_teams.go#L160) | Admin kontrolü |
| POST | `/setup/teams/clear` | [teamsClearHandler](../../web/setup_teams.go#L108) | Admin kontrolü |
| GET | `/setup/teams/generate_wpa_keys` | [teamsGenerateWpaKeysHandler](../../web/setup_teams.go#L231) | Admin kontrolü |
| GET | `/setup/teams/progress` | [teamsUpdateProgressBarHandler](../../web/setup_teams.go#L260) | Doğrudan admin kontrolü yok |
| GET | `/setup/teams/refresh` | [teamsRefreshHandler](../../web/setup_teams.go#L76) | Admin kontrolü |
