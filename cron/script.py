from multiprocessing import Process
import requests
import os

returnCode = 0


def main():
    print("Running Script")
    # START
    # ----------------------------------------------------------------------
    
    response = requests.get(os.getenv("WEBHOOK_URL"))
    if response.status_code == 200:
        print("Webhook Sent Successfully")
        print(response.json())
    else:
        print("Webhook Failed to Send")
        global returnCode
        returnCode = 1
    
    # ----------------------------------------------------------------------
    # END
    print("Script Complete")
    return


if __name__ == "__main__":
    p = Process(target=main)
    p.start()
    p.join()
    exit(returnCode)
