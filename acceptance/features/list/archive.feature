Feature: Create List

  Scenario: OK
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    When I send a "POST" request on "/lists/{{ .ListID }}/archive"
    Then the status code is 200
    And "$.id" is a non empty string
    And "$.name" is a string equal to "Récupération listes de courses"
    And "$.archivationDate" is a non empty string
