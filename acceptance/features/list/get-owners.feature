Feature: Get owner's Lists

  Scenario: OK
    Given I create and authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    When I send a "GET" request on "/lists"
    Then the status code is 200
    And "$.name[0]" is a string equal to "Récupération listes de courses"
    And "$.id[0]" is a non empty string
