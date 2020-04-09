# Read The Room
This demo processes conversations in real-time with the [Amazon Comprehend](https://aws.amazon.com/comprehend/) natural language processing (NLP) service to gain insights about what was said. 

The conversation sentiment, sentiment history, and any detected key words/phrases are displayed in real-time on the dashboard. 

![](https://i.imgur.com/mOapUsp.gif)

## Running The Demo
1. Before you can run this demo you need an AWS account and IAM user credentials. It's important you create and restrict the credentials properly. Learn more about and create the credentials before you get started: [Authentication and Access Control for Amazon Comprehend](https://docs.aws.amazon.com/comprehend/latest/dg/auth-and-access-control.html)

2. Fork/clone this repo.

3. Navigate to the 'application' folder and open '.env'. You may need to show hidden files to see the '.env' file.

4. Add your AWS keys to the '.env' file.
```
AWS_ACCESS_KEY=YOUR_ACCESS_KEY_HERE
AWS_SECRET_KEY=YOUR_SECRET_KEY_HERE
```

5. Use your console to get all the dependencies and run the project (from within the 'application' folder):
`go get -d && go run main.go`

6. A browser window should open automatically and the dashboard should be displayed. Use Google Chrome and navigate to `http://localhost:8091/` if the dashboard does not display or if the dashboard opens in an incompatible browser.

7. Start talking! Try saying positive, neutral, and negative phrases. Does the sentiment match what you're saying? Is it detecting the key words/phrases in your speech?

Try phrases like:

- "I'm so happy I get to play with this Amazon Comprehend Demo."
- "The earth is the third planet from the sun in our solar system."
- "I think eating popcorn with your mouth open is stupid."

## How The Demo Works

1. The [Web Speech API](https://www.google.com/intl/it/chrome/demos/speech.html) is used to convert text to speech.

2. When a chunk of text is ready it is sent to our process api: `http://localhost:8091/process/`

3. The text is forwarded to Amazon Comprehend and the result is returned to the dashboard.

4. [Highcharts](https://www.highcharts.com/) is used to create a visual representation of the sentiment and sentiment history. Any key phrases are formatted and displayed in the text area below the charts. 

## What Else Can You Do With Amazon Comprehend And This Demo

You could use this demo to:

- Monitor in person customer reactions.
- Monitor customer reactions on phone calls and auto recordings.
- Determine how a presentation will be received without an audience.

Some other things you can do with Amazon Comprehend:

- Identify the feature that’s most often mentioned when customers are happy or unhappy about your product.
- Monitor support channels and gain new insights.
- Organize and categorize your documents by topic.


