# **Product Recommendation Microservice**

### **Scenario**:

You are tasked with building a microservice that provides personalized product recommendations to customers based on their previous purchases and browsing behavior.

### **Requirements**:

Implement RESTful API endpoints:
- GET /recommendations?user_id=<id>: returns top 10 recommended product IDs
- POST /activity: accepts user activity like view, purchase, wishlist (mocked payload)
- Simulate user data and product catalog (can be in JSON files or embedded)
- Apply a basic rule-based or collaborative filtering logic (no ML required)
- Store and retrieve recent user interactions from a lightweight in-memory store (e.g., Redis or local map)
- Use Docker to containerize the service
- Include unit and integration tests

### **Bonus Points:**

- Cache recommendations per user
- Implement request/response logging with correlation IDs
