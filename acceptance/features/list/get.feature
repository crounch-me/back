Feature: Get List by id

  @val
  Scenario: OK - List contains products with their categories
    Given I authenticate with a random user
    And I create these lists
      | Name    |
      | Courses |
    And I create these products
      | Name    |
      | Caviar  |
    And I use this body
      """
        {
          "productId": "{{ .ProductID }}",
          "listId": "{{ .ListID }}"
        }
      """
    And I send a "POST" request on "/lists/{{ .ListID }}/products/{{ .ProductID }}"
    And I send a "POST" request on "/lists/{{ .ListID }}/products/40fe3f75-703a-46d8-9520-0d27f7cf4bab"
    When I send a "GET" request on "/lists/{{ .ListID }}"
    Then the status code is 200
    And "$.name" is a string equal to "Courses"
    And "$.id" is a non empty string
    And the returned products from list are
      | ID                                   | Name             | Category name | Buyed |
      | {{ .ProductID }}                     | Caviar           |               | No    |
      | 40fe3f75-703a-46d8-9520-0d27f7cf4bab | Saucisse à cuire | Boucherie     | No    |

  Scenario: KO - User is not the owner
    Given I authenticate with a random user
    And I create these lists
      | name    |
      | Courses |
    And I authenticate with a random user
    When I send a "GET" request on "/lists/{{ .ListID }}"
    Then the status code is 403
    And "$.error" is a string equal to "forbidden-error"

  Scenario: KO - List id not UUID
    Given I authenticate with a random user
    When I send a "GET" request on "/lists/a"
    Then the status code is 400
    And "$.fields[0].name" is a string equal to "listID"
    And "$.fields[0].error" is a string equal to "uuid"

  Scenario: KO - List not found
    Given I authenticate with a random user
    When I send a "GET" request on "/lists/00000000-0000-0000-0000-000000000000"
    Then the status code is 404
    And "$.error" is a string equal to "list-not-found-error"
