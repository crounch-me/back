Feature: Get Open Food Facts Product from a List

  Scenario: OK
    Given I create and authenticate with a random user
    And I create a random list
    And I add a random off product
    When I send a "GET" request on "/lists/{{ .ListID }}/offproducts"
    Then the status code is 200
    And "$.code[0]" is a string equal to "{{ .OFFProductCode }}"
