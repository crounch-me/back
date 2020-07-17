Feature: Delete List

  Scenario: OK
    Given I authenticate with a random user
    And I create these lists
      | name                           |
      | Récupération listes de courses |
    When I send a "DELETE" request on "/lists/{{ .ListID }}"
    Then the status code is 204

  Scenario: KO - List id is not an UUID
    Given I authenticate with a random user
    When I send a "DELETE" request on "/lists/a"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    And "$.fields[0].name" has string value "listID"
    And "$.fields[0].error" has string value "uuid"

  Scenario: KO - List not found
    Given I authenticate with a random user
    When I send a "DELETE" request on "/lists/00000000-0000-0000-0000-000000000000"
    Then the status code is 404
    And "$.error" has string value "list-not-found-error"
