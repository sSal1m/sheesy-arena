# Veri yapıları ve alan sözlüğü

`model/` yapıları ile temel oyun yapılarının kaynak kodundan çıkarılmış alanlarıdır. Her yapı bir Bolt tablosu değildir; gerçek 13 bucket ana haritada listelenir. Tipler Go gösterimiyle korunur. `db:"id,manual"` uygulamanın verdiği anahtarı; `db:"id"` otomatik artan anahtarı gösterir. Canlı veri, parola veya API anahtarı değerleri bu dosyaya okunmaz.

## Foul

[game/foul.go](../../game/foul.go#L8)

| Alan | Tip / etiket |
| --- | --- |
| `FoulId` | `int` |
| `IsMajor` | `bool` |
| `TeamId` | `int` |
| `RuleId` | `int` |

## Hub

[game/hub.go](../../game/hub.go#L12)

| Alan | Tip / etiket |
| --- | --- |
| `WonAuto` | `bool` |
| `ShiftCounts` | `[ShiftCount]int` |

## RankingFields

[game/ranking_fields.go](../../game/ranking_fields.go#L10)

| Alan | Tip / etiket |
| --- | --- |
| `RankingPoints` | `int` |
| `MatchPoints` | `int` |
| `AutoFuelPoints` | `int` |
| `TowerPoints` | `int` |
| `Random` | `float64` |
| `Wins` | `int` |
| `Losses` | `int` |
| `Ties` | `int` |
| `Disqualifications` | `int` |
| `Played` | `int` |

## Ranking

[game/ranking_fields.go](../../game/ranking_fields.go#L23)

| Alan | Tip / etiket |
| --- | --- |
| `TeamId` | `` int `db:"id,manual"` `` |
| `Rank` | `int` |
| `PreviousRank` | `int` |
| `(gömülü)` | `RankingFields` |

## Score

[game/score.go](../../game/score.go#L8)

| Alan | Tip / etiket |
| --- | --- |
| `AutoTowerStatuses` | `[3]TowerStatus` |
| `Hub` | `Hub` |
| `EndgameTowerStatuses` | `[3]TowerStatus` |
| `Fouls` | `[]Foul` |
| `PlayoffDq` | `bool` |

## ScoreSummary

[game/score_summary.go](../../game/score_summary.go#L8)

| Alan | Tip / etiket |
| --- | --- |
| `AutoFuelPoints` | `int` |
| `AutoTowerPoints` | `int` |
| `TeleopFuelPoints` | `int` |
| `TeleopTowerPoints` | `int` |
| `NumFuel` | `int` |
| `NumFuelPostMatch` | `int` |
| `NumFuelGoal` | `int` |
| `MatchPoints` | `int` |
| `PostMatchPoints` | `int` |
| `FoulPoints` | `int` |
| `Score` | `int` |
| `PlayoffDq` | `bool` |
| `EnergizedBonusRankingPoint` | `bool` |
| `SuperchargedBonusRankingPoint` | `bool` |
| `TraversalBonusRankingPoint` | `bool` |
| `BonusRankingPoints` | `int` |
| `NumOpponentMajorFouls` | `int` |

## Alliance

[model/alliance.go](../../model/alliance.go#L10)

| Alan | Tip / etiket |
| --- | --- |
| `Id` | `` int `db:"id,manual"` `` |
| `TeamIds` | `[]int` |
| `Lineup` | `[3]int` |
| `YellowCard` | `bool` |

## AllianceSelectionRankedTeam

[model/alliance.go](../../model/alliance.go#L17)

| Alan | Tip / etiket |
| --- | --- |
| `Rank` | `int` |
| `TeamId` | `int` |
| `Picked` | `bool` |

## Award

[model/award.go](../../model/award.go#L13)

| Alan | Tip / etiket |
| --- | --- |
| `Id` | `` int `db:"id"` `` |
| `Type` | `AwardType` |
| `AwardName` | `string` |
| `TeamId` | `int` |
| `PersonName` | `string` |
| `TbaTypeEnum` | `*int` |

## Database

[model/database.go](../../model/database.go#L23)

| Alan | Tip / etiket |
| --- | --- |
| `Path` | `string` |
| `bolt` | `*bbolt.DB` |
| `allianceTable` | `*table[Alliance]` |
| `awardTable` | `*table[Award]` |
| `eventSettingsTable` | `*table[EventSettings]` |
| `judgingSlotTable` | `*table[JudgingSlot]` |
| `lowerThirdTable` | `*table[LowerThird]` |
| `matchTable` | `*table[Match]` |
| `matchResultTable` | `*table[MatchResult]` |
| `rankingTable` | `*table[game.Ranking]` |
| `scheduleBlockTable` | `*table[ScheduleBlock]` |
| `scheduledBreakTable` | `*table[ScheduledBreak]` |
| `sponsorSlideTable` | `*table[SponsorSlide]` |
| `teamTable` | `*table[Team]` |
| `userSessionTable` | `*table[UserSession]` |

## EventSettings

[model/event_settings.go](../../model/event_settings.go#L41)

| Alan | Tip / etiket |
| --- | --- |
| `Id` | `` int `db:"id"` `` |
| `Name` | `string` |
| `PlayoffType` | `PlayoffType` |
| `NumPlayoffAlliances` | `int` |
| `SelectionRound2Order` | `string` |
| `SelectionRound3Order` | `string` |
| `SelectionShowUnpickedTeams` | `bool` |
| `TbaDownloadEnabled` | `bool` |
| `TbaPublishingEnabled` | `bool` |
| `TbaEventCode` | `string` |
| `TbaSecretId` | `string` |
| `TbaSecret` | `string` |
| `AutoAudienceDisplayEnabled` | `bool` |
| `NexusEnabled` | `bool` |
| `NexusAutoQueueEnabled` | `bool` |
| `NexusAutoQueueKey` | `string` |
| `NetworkSecurityEnabled` | `bool` |
| `ApAddress` | `string` |
| `ApPassword` | `string` |
| `ApChannel` | `int` |
| `SwitchAddress` | `string` |
| `SwitchPassword` | `string` |
| `SCCManagementEnabled` | `bool` |
| `RedSCCAddress` | `string` |
| `BlueSCCAddress` | `string` |
| `SCCUsername` | `string` |
| `SCCPassword` | `string` |
| `SCCUpCommands` | `string` |
| `SCCDownCommands` | `string` |
| `PlcAddress` | `string` |
| `LedControllerAddress` | `string` |
| `LedUniverseMode` | `string` |
| `AdminPassword` | `string` |
| `TeamSignRed1Id` | `int` |
| `TeamSignRed2Id` | `int` |
| `TeamSignRed3Id` | `int` |
| `TeamSignRedTimerId` | `int` |
| `TeamSignBlue1Id` | `int` |
| `TeamSignBlue2Id` | `int` |
| `TeamSignBlue3Id` | `int` |
| `TeamSignBlueTimerId` | `int` |
| `UseLiteUdpPort` | `bool` |
| `BlackmagicAddresses` | `string` |
| `CompanionAddress` | `string` |
| `CompanionPort` | `int` |
| `CompanionMatchPreviewPage` | `int` |
| `CompanionMatchPreviewRow` | `int` |
| `CompanionMatchPreviewColumn` | `int` |
| `CompanionSetAudiencePage` | `int` |
| `CompanionSetAudienceRow` | `int` |
| `CompanionSetAudienceColumn` | `int` |
| `CompanionMatchStartPage` | `int` |
| `CompanionMatchStartRow` | `int` |
| `CompanionMatchStartColumn` | `int` |
| `CompanionTeleopStartPage` | `int` |
| `CompanionTeleopStartRow` | `int` |
| `CompanionTeleopStartColumn` | `int` |
| `CompanionEndgameStartPage` | `int` |
| `CompanionEndgameStartRow` | `int` |
| `CompanionEndgameStartColumn` | `int` |
| `CompanionMatchEndPage` | `int` |
| `CompanionMatchEndRow` | `int` |
| `CompanionMatchEndColumn` | `int` |
| `CompanionPostResultPage` | `int` |
| `CompanionPostResultRow` | `int` |
| `CompanionPostResultColumn` | `int` |
| `CompanionAllianceSelectionPage` | `int` |
| `CompanionAllianceSelectionRow` | `int` |
| `CompanionAllianceSelectionColumn` | `int` |
| `CompanionMatchAbortPage` | `int` |
| `CompanionMatchAbortRow` | `int` |
| `CompanionMatchAbortColumn` | `int` |
| `AutoDurationSec` | `int` |
| `PauseDurationSec` | `int` |
| `TransitionShiftDurationSec` | `int` |
| `ShiftDurationSec` | `int` |
| `EndgameDurationSec` | `int` |
| `EnergizedBonusThreshold` | `int` |
| `SuperchargedBonusThreshold` | `int` |
| `TraversalBonusThreshold` | `int` |

## JudgingSlot

[model/judging_slot.go](../../model/judging_slot.go#L13)

| Alan | Tip / etiket |
| --- | --- |
| `Id` | `` int `db:"id"` `` |
| `Time` | `time.Time` |
| `TeamId` | `int` |
| `PreviousMatchNumber` | `int` |
| `PreviousMatchTime` | `time.Time` |
| `NextMatchNumber` | `int` |
| `NextMatchTime` | `time.Time` |
| `JudgeNumber` | `int` |

## LowerThird

[model/lower_third.go](../../model/lower_third.go#L12)

| Alan | Tip / etiket |
| --- | --- |
| `Id` | `` int `db:"id"` `` |
| `TopText` | `string` |
| `BottomText` | `string` |
| `DisplayOrder` | `int` |
| `AwardId` | `int` |

## Match

[model/match.go](../../model/match.go#L31)

| Alan | Tip / etiket |
| --- | --- |
| `Id` | `` int `db:"id"` `` |
| `Type` | `MatchType` |
| `TypeOrder` | `int` |
| `Time` | `time.Time` |
| `LongName` | `string` |
| `ShortName` | `string` |
| `NameDetail` | `string` |
| `PlayoffMatchGroupId` | `string` |
| `PlayoffRedAlliance` | `int` |
| `PlayoffBlueAlliance` | `int` |
| `Red1` | `int` |
| `Red1IsSurrogate` | `bool` |
| `Red2` | `int` |
| `Red2IsSurrogate` | `bool` |
| `Red3` | `int` |
| `Red3IsSurrogate` | `bool` |
| `Blue1` | `int` |
| `Blue1IsSurrogate` | `bool` |
| `Blue2` | `int` |
| `Blue2IsSurrogate` | `bool` |
| `Blue3` | `int` |
| `Blue3IsSurrogate` | `bool` |
| `StartedAt` | `time.Time` |
| `ScoreCommittedAt` | `time.Time` |
| `FieldReadyAt` | `time.Time` |
| `Status` | `game.MatchStatus` |
| `UseTiebreakCriteria` | `bool` |
| `TbaMatchKey` | `TbaMatchKey` |

## TbaMatchKey

[model/match.go](../../model/match.go#L62)

| Alan | Tip / etiket |
| --- | --- |
| `CompLevel` | `string` |
| `SetNumber` | `int` |
| `MatchNumber` | `int` |

## MatchResult

[model/match_result.go](../../model/match_result.go#L12)

| Alan | Tip / etiket |
| --- | --- |
| `Id` | `` int `db:"id"` `` |
| `MatchId` | `int` |
| `PlayNumber` | `int` |
| `MatchType` | `MatchType` |
| `RedScore` | `*game.Score` |
| `BlueScore` | `*game.Score` |
| `RedCards` | `map[string]string` |
| `BlueCards` | `map[string]string` |
| `PlayoffRedAllianceCard` | `string` |
| `PlayoffBlueAllianceCard` | `string` |

## ScheduleBlock

[model/schedule_block.go](../../model/schedule_block.go#L13)

| Alan | Tip / etiket |
| --- | --- |
| `Id` | `` int `db:"id"` `` |
| `MatchType` | `MatchType` |
| `StartTime` | `time.Time` |
| `NumMatches` | `int` |
| `MatchSpacingSec` | `int` |

## ScheduledBreak

[model/scheduled_break.go](../../model/scheduled_break.go#L13)

| Alan | Tip / etiket |
| --- | --- |
| `Id` | `` int `db:"id"` `` |
| `MatchType` | `MatchType` |
| `TypeOrderBefore` | `int` |
| `Time` | `time.Time` |
| `DurationSec` | `int` |
| `Description` | `string` |

## SponsorSlide

[model/sponsor_slide.go](../../model/sponsor_slide.go#L10)

| Alan | Tip / etiket |
| --- | --- |
| `Id` | `` int `db:"id"` `` |
| `Subtitle` | `string` |
| `Line1` | `string` |
| `Line2` | `string` |
| `Image` | `string` |
| `DisplayTimeSec` | `int` |
| `DisplayOrder` | `int` |

## Team

[model/team.go](../../model/team.go#L10)

| Alan | Tip / etiket |
| --- | --- |
| `Id` | `` int `db:"id,manual"` `` |
| `Name` | `string` |
| `Nickname` | `string` |
| `City` | `string` |
| `StateProv` | `string` |
| `Country` | `string` |
| `SchoolName` | `string` |
| `RookieYear` | `int` |
| `RobotName` | `string` |
| `Accomplishments` | `string` |
| `WpaKey` | `string` |
| `YellowCard` | `bool` |
| `HasConnected` | `bool` |
| `FtaNotes` | `string` |

## UserSession

[model/user_session.go](../../model/user_session.go#L10)

| Alan | Tip / etiket |
| --- | --- |
| `Id` | `` int `db:"id"` `` |
| `Token` | `string` |
| `Username` | `string` |
| `CreatedAt` | `time.Time` |
