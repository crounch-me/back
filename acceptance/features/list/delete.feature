Feature: Delete List

  Scenario: OK
    Given I create and authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    When I send a "DELETE" request on "/lists/{{ .ListID }}"
    Then the status code is 204
