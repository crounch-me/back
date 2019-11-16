Feature: List

  Scenario: Create a list
    Given I create a random user
    And I'm authenticated with this random user
    And I use this body
      """
        {
          "name": "Ma liste de courses"
        }
      """
    When I send a "POST" request on "/lists"
    Then the status code is 201
    And "$.name" is a string equal to "Ma liste de courses"
    And "$.id" is a non empty string
