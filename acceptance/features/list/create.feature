Feature: Create List

  Scenario: Missing list name
    Given I create and authenticate with a random user
    And I use this body
      """
        {}
      """
    When I send a "POST" request on "/lists"
    Then the status code is 400
    And "$.key[0]" is a string equal to "name"
    And "$.tag[0]" is a string equal to "required"

  Scenario: OK
    Given I create and authenticate with a random user
    And I use this body
      """
        {
          "name": "Courses"
        }
      """
    When I send a "POST" request on "/lists"
    Then the status code is 201
    And "$.name" is a string equal to "Courses"
    And "$.id" is a non empty string
