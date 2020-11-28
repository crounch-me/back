package list

import "github.com/crounch-me/back/internal"

func (ls *ListService) isUserAuthorized(listID string, userID string) (bool, *internal.Error) {
	contributorIDs, err := ls.ContributorStorage.GetContributorsIDs(listID)
	if err != nil {
		return false, err
	}

	for _, contributorID := range contributorIDs {
		if contributorID == userID {
			return true, nil
		}
	}

	return false, nil
}
