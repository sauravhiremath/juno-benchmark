#include <stdio.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <string.h>
#include <arpa/inet.h>
#include <stdlib.h>
#include <fcntl.h>  // for open
#include <unistd.h> // for close
#include <pthread.h>
#include <errno.h>

#define BUFFER_SIZE 256
char buf[BUFFER_SIZE] = {0};
size_t buf_idx = 0;

void *clientThread(void *arg)
{
    printf("In thread\n");
    char message[2000];
    char buffer[4096];
    int clientSocket;
    ssize_t numBytesSent = 0;
    struct sockaddr_in serverAddr;
    socklen_t addr_size;
    // Create the socket.
    clientSocket = socket(PF_INET, SOCK_STREAM, 0);
    //Configure settings of the server address
    // Address family is Internet
    serverAddr.sin_family = AF_INET;
    //Set port number, using htons function
    serverAddr.sin_port = htons(4000);
    //Set IP address to localhost
    serverAddr.sin_addr.s_addr = inet_addr("127.0.0.1");
    memset(serverAddr.sin_zero, '\0', sizeof serverAddr.sin_zero);
    //Connect the socket to the server using the address
    addr_size = sizeof serverAddr;
    connect(clientSocket, (struct sockaddr *)&serverAddr, addr_size);
    printf("connected\n");
    // Send initialisation request
    strcpy(message, "{\"requestId\":\"module1-1xsx\",\"type\":1,\"moduleId\":\"module11\",\"version\":\"1.0.0\"}\n");
    numBytesSent = send(clientSocket, message, strlen(message), 0);
    if (numBytesSent < 0)
    {
        printf("Send failed\n");
    }
    printf("sent message size %d \n", numBytesSent);

    //Read the message from the server into the buffer
    while (buf_idx < BUFFER_SIZE && 1 == read(clientSocket, &buf[buf_idx], 1))
    {
        if (buf_idx > 0 &&
            '\n' == buf[buf_idx])
        {
            break;
        }
        buf_idx++;
    }

    //Print the received message
    printf("Data received: %s\n", buf);
    close(clientSocket);
    pthread_exit(NULL);
}

int main()
{
    int i = 0;
    pthread_t tid[51];
    while (i < 50)
    {
        if (pthread_create(&tid[i], NULL, clientThread, NULL) != 0)
            printf("Failed to create thread\n");
        i++;
    }
    sleep(1);
    i = 0;
    while (i < 1)
    {
        pthread_join(tid[i++], NULL);
        printf("%d:\n", i);
    }
    return 0;
}