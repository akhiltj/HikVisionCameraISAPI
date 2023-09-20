package hikvision

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	// "os"
	"strings"
	"github.com/icholy/digest"
)

type HttpEventReader struct {
	Debug  bool
	client *http.Client
}

type CustomEvent struct {
	IPAddress        string `xml:"ipAddress"`
	IPv6Address      string `xml:"ipv6Address"`
	PortNo           int    `xml:"portNo"`
	Protocol         string `xml:"protocol"`
	MACAddress       string `xml:"macAddress"`
	ChannelID        int    `xml:"channelID"`
	DateTime         string `xml:"dateTime"`
	ActivePostCount  int    `xml:"activePostCount"`
	EventType        string `xml:"eventType"`
	EventState       string `xml:"eventState"`
	EventDescription string `xml:"eventDescription"`
	ChannelName      string `xml:"channelName"`
	Camera           *HikCamera
}

func (eventReader *HttpEventReader) ReadEvents(camera *HikCamera, channel chan<- HikEvent, callback func(),sendmqttmessage func(data TestMessage)) {
	if eventReader.client == nil {
		eventReader.client = &http.Client{}
		if camera.AuthMethod == Digest {
			eventReader.client.Transport = &digest.Transport{
				Username: camera.Username,
				Password: camera.Password,
			}
		}
	}

	request, err := http.NewRequest("GET", camera.Url+"Event/notification/alertStream", nil)
	if err != nil {
		fmt.Printf("HIK: Error: Could not connect to camera %s\n", camera.Name)
		fmt.Println("HIK: Error", err)
		callback()
		return
	}
	if camera.AuthMethod == Basic {
		request.SetBasicAuth(camera.Username, camera.Password)
	}

	response, err := eventReader.client.Do(request)
	if err != nil {
		fmt.Printf("HIK: Error opening HTTP connection to camera %s\n", camera.Name)
		fmt.Println(err)
		return
	}

	if response.StatusCode != 200 {
		fmt.Printf("HIK: BAD STATUS %d", response.StatusCode)
	}
	defer response.Body.Close()

	// FIGURE OUT MULTIPART BOUNDARY
	mediaType, params, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if mediaType != "multipart/mixed" || params["boundary"] == "" {
		fmt.Println("HIK: ERROR: Camera " + camera.Name + " does not seem to support event streaming")
		fmt.Println("            Is it a doorbell? Try adding rawTcp to its config!")
		callback()
		return
	}
	multipartBoundary := params["boundary"]

	customEvent := CustomEvent{}

	// READ PART BY PART
	multipartReader := multipart.NewReader(response.Body, multipartBoundary)
	// fmt.Println("------response body start-------")
	// _, err = io.Copy(os.Stdout, response.Body) // Print the response body to stdout
	// fmt.Println("\n------response body end---------")

	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			// End of the multipart content
			break
		}
		if err != nil {
			// fmt.Println("Error reading part:", err)
			continue
		}

		// Read the entire part body
		body, err := io.ReadAll(part)
		if err != nil {
			fmt.Println("Error reading part body:", err)
			continue
		}

		// Create a reader for the body content
		bodyReader := strings.NewReader(string(body))

		// Create a new decoder for the XML content
		decoder := xml.NewDecoder(bodyReader)

		// Decode the XML into the CustomEvent struct
		if err := decoder.Decode(&customEvent); err != nil {
			fmt.Println("Error decoding XML:", err)
			continue
		}

		// Check if the event type is "VMD" and the event state is "active"
		if customEvent.EventType == "VMD" && customEvent.EventState == "active" {
			// This is a motion alarm event, you can handle it here
			fmt.Printf("Handling Motion Alarm Event\n")

			// Rest of your code for handling the motion detection event
			// ...

			// FILL IN THE CAMERA INTO FRESHLY-UNMARSHALLED EVENT
			customEvent.Camera = camera

			if eventReader.Debug {
				log.Printf("%s event: %s (%s - %d)", customEvent.Camera.Name, customEvent.EventType, customEvent.EventState, customEvent.ActivePostCount)
			}

			switch customEvent.EventState {
			case "active":
				// messageHandler(camera.Name, customEvent.EventType, "Motion alarm detected")
				fmt.Printf("Motion alarm Event Active\n")
				// Handle active state
			case "inactive":
				// Handle inactive state
			}
		} else if customEvent.EventType == "videoloss" && customEvent.EventState == "inactive" {
			// This is a video loss event, you can handle it here
			fmt.Printf("Handling Video Loss Event\n")
			data := TestMessage{
				Type: "42",
				Message: "videoloss",
			}
			sendmqttmessage(data)

			// Rest of your code for handling the video loss event
			// ...

			// FILL IN THE CAMERA INTO FRESHLY-UNMARSHALLED EVENT
			customEvent.Camera = camera

			if eventReader.Debug {
				log.Printf("%s event: %s (%s - %d)", customEvent.Camera.Name, customEvent.EventType, customEvent.EventState, customEvent.ActivePostCount)
			}

			switch customEvent.EventState {
			case "active":
				fmt.Printf("video loss alarm Event Active\n")
				// Handle active state
			case "inactive":
				// fmt.Printf("%s event: %s (%s - %d)", customEvent.Camera.Name, customEvent.EventType, customEvent.EventState, customEvent.ActivePostCount)
				// Handle inactive state
			}
		} else if customEvent.EventType == "regionEntrance" && customEvent.EventState == "inactive" {
			// This is a video loss event, you can handle it here
			fmt.Printf("regionEntrance Event\n")

			// Rest of your code for handling the video loss event
			// ...

			// FILL IN THE CAMERA INTO FRESHLY-UNMARSHALLED EVENT
			customEvent.Camera = camera

			if eventReader.Debug {
				log.Printf("%s event: %s (%s - %d)", customEvent.Camera.Name, customEvent.EventType, customEvent.EventState, customEvent.ActivePostCount)
			}

			switch customEvent.EventState {
			case "active":
				fmt.Printf("regionEntrance Active\n")
				//mqtt.SendMessage(config.TopicRoot+"/alarmserver", `{ "status": "up" }`)
				// messageHandler(camera.Name, customEvent.EventType, "video loss alarm detected")
				// Handle active state
			case "inactive":
				fmt.Printf("regionEntrance Inactive\n")
				// Handle inactive state
			}
		} else if customEvent.EventType == "linedetection" && customEvent.EventState == "inactive" {
			// This is a video loss event, you can handle it here
			fmt.Printf("linedetection Event\n")

			// Rest of your code for handling the video loss event
			// ...

			// FILL IN THE CAMERA INTO FRESHLY-UNMARSHALLED EVENT
			customEvent.Camera = camera

			if eventReader.Debug {
				log.Printf("%s event: %s (%s - %d)", customEvent.Camera.Name, customEvent.EventType, customEvent.EventState, customEvent.ActivePostCount)
			}

			switch customEvent.EventState {
			case "active":
				fmt.Printf("linedetection Active\n")
				//mqtt.SendMessage(config.TopicRoot+"/alarmserver", `{ "status": "up" }`)
				// messageHandler(camera.Name, customEvent.EventType, "video loss alarm detected")
				// Handle active state
			case "inactive":
				fmt.Printf("linedetection Inactive\n")
				// Handle inactive state
			}
		}
	}
}


