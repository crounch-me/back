Feature: Get owner's Lists

  Scenario: OK
    Given I authenticate with a random user
    And I create these lists
      | name    |
      | Courses |
    When I send a "GET" request on "/lists"
    Then the status code is 200
    And "$.name[0]" is a string equal to "Courses"
    And "$.id[0]" is a non empty string

  Scenario: OK - Doesn't get other users lists
    Given I authenticate with a random user
    And I create these lists
      | name    |
      | Courses |
    And I authenticate with a random user
    When I send a "GET" request on "/lists"
    Then the status code is 200
    And the body is an empty array
