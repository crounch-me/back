# Feature: Search Default Products by Name

#   Scenario: OK - Beginning of the name
#     Given I authenticate with a random user
#     And I use this body
#       """
#         {
#           "name": "len"
#         }
#       """
#     When I send a "POST" request on "/products/search"
#     Then the status code is 200
#     And "$[0].id" is a non empty string
#     And "$[0].name" is a string equal to "Lentille"
#     And "$[0].category.name" is a string equal to "Epicerie"

#   Scenario: OK - Middle of the name
#     Given I authenticate with a random user
#     And I use this body
#       """
#         {
#           "name": "mor"
#         }
#       """
#     When I send a "POST" request on "/products/search"
#     Then the status code is 200
#     And "$[0].id" is a non empty string
#     And "$[0].name" is a string equal to "Saucisse de Morteau"
#     And "$[0].category.name" is a string equal to "Boucherie"

#   Scenario: OK - Accentuated character
#     Given I authenticate with a random user
#     And I use this body
#       """
#         {
#           "name": "Montbeliard"
#         }
#       """
#     When I send a "POST" request on "/products/search"
#     Then the status code is 200
#     And "$[0].id" is a non empty string
#     And "$[0].name" is a string equal to "Saucisse de Montbéliard"
#     And "$[0].category.name" is a string equal to "Boucherie"

#   Scenario: Invalid body
#     Given I authenticate with a random user
#     And I use an invalid body
#     When I send a "POST" request on "/products/search"
#     Then the status code is 400
#     And "$.error" has string value "unmarshal-error"

#   Scenario: KO - Missing fields
#     Given I use an empty valid body
#     And I send a "POST" request on "/products/search"
#     Then the status code is 400
#     And "$.error" has string value "invalid-error"
#     And "$.fields[0].name" has string value "name"
#     And "$.fields[0].error" has string value "gt"

#   Scenario: KO - Too short name
#     Given I authenticate with a random user
#     And I use this body
#       """
#         {
#           "name": "mo"
#         }
#       """
#     When I send a "POST" request on "/products/search"
#     Then the status code is 400
#     And "$.error" has string value "invalid-error"
#     And "$.fields[0].name" has string value "name"
#     And "$.fields[0].error" has string value "gt"

#   Scenario: KO - Too long name
#     Given I authenticate with a random user
#     And I use this body
#       """
#         {
#           "name": "mooooooooooooooo"
#         }
#       """
#     When I send a "POST" request on "/products/search"
#     Then the status code is 400
#     And "$.error" has string value "invalid-error"
#     And "$.fields[0].name" has string value "name"
#     And "$.fields[0].error" has string value "lt"
