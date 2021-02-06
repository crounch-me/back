@listing
Feature: Create List

  Scenario: OK
    Given I authenticate with a random user
    And I use this body
      """
        {
          "name": "Courses"
        }
      """
    When I send a "POST" request on "/listing/lists"
    Then the status code is 201

  Scenario: KO - Invalid body
    Given I authenticate with a random user
    And I use an invalid body
    When I send a "POST" request on "/listing/lists"
    Then the status code is 400
    And "$.error" has string value "unmarshal-error"

  Scenario: KO - Missing fields
    Given I authenticate with a random user
    And I use an empty valid body
    And I send a "POST" request on "/listing/lists"
    Then the status code is 400
    And "$.error" has string value "invalid-error"
    And "$.fields[0].name" has string value "name"
    And "$.fields[0].error" has string value "required"
