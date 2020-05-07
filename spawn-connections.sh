# Install tcpkali (present in AUR library)

# Initialization
tcpkali -em "{\"requestId\":\"module1-1234567890\",\"type\":1,\"moduleId\":\"module11\",\"version\":\"1.0.0\"}\n" 127.0.0.1:4000 -T 1000 -c 32
# Function call
tcpkali -em "{\r\n    \"requestId\": \"module1-1234567890\",\r\n    \"type\": 3,\r\n    \"function\": \"module2.calculateSum\",\r\n    \"data\": {\r\n        \"values\": [\r\n            1,\r\n            2,\r\n            3,\r\n            4,\r\n            5\r\n        ]\r\n    }\r\n}\n" 127.0.0.1:4000 -T 1000 -c 32

# Hook register
tcpkali -em "{\r\n    \"requestId\": \"module1-1234567890\",\r\n    \"type\": 5,\r\n    \"hook\": \"users.passwordChanged\" \/\/ Will listen for the \"passwordChanged\" hook of the \"users\" module\r\n}\n" 127.0.0.1:4000 -T 1000 -c 32
# Trigger hook
tcpkali -em "{\r\n    \"requestId\": \"module1-1234567890\",\r\n    \"type\": 7,\r\n    \"hook\": \"passwordChanged\",\r\n    \"data\": {\r\n        \"userId\": \"testUser\"\r\n    }\r\n}\n" 127.0.0.1:4000 -T 1000 -c 32

# Declare function
tcpkali -em "{\r\n    \"requestId\": \"module1-1234567890\",\r\n    \"type\": 9,\r\n    \"function\": \"passwordChanged\"\r\n}\n" 127.0.0.1:4000 -T 1000 -c 32



