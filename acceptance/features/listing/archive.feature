@listing
Feature: Create List

  Scenario: OK
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    When I send a "POST" request on "/listing/lists/{{ .ListID }}/archive"
    Then the status code is 204
