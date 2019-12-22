Feature: Add Open Food Facts Product to List

  @current
  Scenario: OK
    Given I create and authenticate with a random user
    And I create a random list
    And I use this body
      """
        {
          "code": "1234567890123"
        }
      """
    When I send a "POST" request on "/lists/{{ .ListID }}/offproducts"
    Then the status code is 201
    And "$.code" is a string equal to "1234567890123"
    And "$.id" is a non empty string
