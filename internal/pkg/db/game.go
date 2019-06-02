package db

var DefaultScoreInc = 25
var DefaultScoreDec = 25

func UpdateScoresAndWinRate(winnerNick string, loserNick string, scoreInc int, scoreDec int) error {

	row, err := findRowBy(AuthDbName, AuthUsersTable, "nickname", "id=?", winnerNick)
	if err != nil {
		return err
	}
	var winnerId int
	err = row.Scan(&winnerId)
	if err != nil {
		return err
	}

	row, err = findRowBy(AuthDbName, AuthUsersTable, "nickname", "id=?", loserNick)
	if err != nil {
		return err
	}
	var loserId int
	err = row.Scan(&loserId)
	if err != nil {
		return err
	}

	_, err = updateBy(CoreDbName, CoreUsersTable, "score=score+?, win=win+1, games=games+1", "id=?", scoreInc, winnerId)
	if err != nil {
		return err
	}

	_, err = updateBy(CoreDbName, CoreUsersTable, "score=score-?, lose=lose+1, games=games+1", "id=?", scoreDec, loserId)
	if err != nil {
		return err
	}

	return nil
}
