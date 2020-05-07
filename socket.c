#include <stdio.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <string.h>
#include <arpa/inet.h>
#include <stdlib.h>
#include <fcntl.h>  // for open
#include <unistd.h> // for close
#include <pthread.h>
#include<time.h>

// #################################################
#define BUFFER_SIZE 256            // MAX len of message to send and recieve from juno (in Bytes)
#define MAX_THREADS 1000          // MAX threads to spawn
#define MAX_CONNECTIONS 1000      // Max connections per thread
// #################################################

char buf[BUFFER_SIZE] = {0};
char send_buf[BUFFER_SIZE] = {0};
size_t buf_idx = 0;
size_t movAvg = 0;

void calcMovingAvg(double);

void *clientThread(void *arg)
{
    clock_t t;
    int clientSocket;
    ssize_t numBytesSent = 0;
    struct sockaddr_in serverAddr;
    socklen_t addr_size;
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
    int connections = MAX_CONNECTIONS;
    // #################################################################################
    while(connections > 0) {
        t = clock();
        // Create the socket.
        clientSocket = socket(PF_INET, SOCK_STREAM, 0);
        connect(clientSocket, (struct sockaddr *)&serverAddr, addr_size);
        // Send initialisation request
        strcpy(send_buf, "{\"requestId\":\"module1-1234567890\",\"type\":1,\"moduleId\":\"module11\",\"version\":\"1.0.0\"}\n");
        numBytesSent = send(clientSocket, send_buf, strlen(send_buf), 0);
        if (numBytesSent < 0)
        {
            printf("Send failed\n");
        }
        // printf("sent message size %d \n", numBytesSent);

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
        // printf("Data received: %s\n", buf);
        t = clock() - t;
        close(clientSocket);
        connections--;
        memset(buf, 0, sizeof(buf));
        memset(send_buf, 0, sizeof(send_buf));
        buf_idx = 0;
        numBytesSent = 0;

        double time_taken = ((double)t)/CLOCKS_PER_SEC;
        movAvg = 0.8*time_taken + 0.2*movAvg;
        printf("Current time: %f secs \nMoving Average: %f secs\n\n", time_taken, movAvg);
    }
    // ##################################################################################
    pthread_exit(NULL);
}

int main()
{
    int i = 0;
    pthread_t tid[MAX_THREADS];
    while (i < MAX_THREADS-1)
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