package lists

import "github.com/crounch-me/back/domain"

func (ls *ListService) isUserAuthorized(listID string, userID string) (bool, *domain.Error) {
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
