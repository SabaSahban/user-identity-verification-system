# Web-Based Identity Verification System

This project implements a web-based identity verification service using two microservices: one for handling user request submissions and storing data securely, and another for verifying user identity through facial recognition APIs. The system uses RabbitMQ for task queuing and a web email service to notify users of their verification status.

## Features
- **User Request Submission:** Collects user details (email, national ID, IP address, and photos) and stores them in a web database.
- **Identity Verification:** Uses facial recognition and similarity APIs to compare user-submitted images.
- **Task Queueing:** Utilizes RabbitMQ for task management between microservices.
- **Email Notifications:** Sends status updates to users about their identity verification progress.

## Architecture
<img width="630" alt="Screenshot 1403-06-19 at 18 08 36" src="https://github.com/user-attachments/assets/d3c5fe9b-8f3c-4239-a47c-0e72ec149689">


## Services Used
- **Web Database:** Stores user information securely, with encrypted fields for sensitive data like the national ID.
- **Object Storage:** Stores user-submitted photos and enables retrieval for identity verification.
- **RabbitMQ:** Manages task queues for processing identity verification requests.
- **Facial Recognition & Similarity APIs:** Processes user images to verify identity.
- **Email Service:** Sends confirmation or rejection emails after identity verification.

## How It Works
1. **User Request Submission:**  
   - The user submits their information via an API, including two images of themselves.
   - The system stores the data and pushes the username to a RabbitMQ queue.
  

2. **Identity Verification:**  
   - The second microservice processes the RabbitMQ queue and verifies the images using facial recognition and similarity checks.
   - Depending on the result, it updates the database and notifies the user via email.


## Installation & Setup
1. Clone the repository.
2. Set up your web services for database, object storage, RabbitMQ, facial recognition, and email notifications.
3. Configure environment variables for your web service credentials.
4. Deploy the microservices using Docker or your preferred deployment method.
