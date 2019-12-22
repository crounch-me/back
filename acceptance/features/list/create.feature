Feature: Create List

  Scenario: OK
    Given I create and authenticate with a random user
    And I use this body
      """
        {
          "name": "Creation liste de courses"
        }
      """
    When I send a "POST" request on "/lists"
    Then the status code is 201
    And "$.name" is a string equal to "Creation liste de courses"
    And "$.id" is a non empty string
