Feature: Create List

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
