---
title: "Ensign for Data Scientists"
weight: 50
bookFlatSection: false
bookToc: true
bookHidden: false
bookCollapseSection: false
bookSearchExclude: false
---

# Ensign for Data Scientists

Here's a question we frequently get from our data scientist friends:
> What does event-driven data science even look like??

In this tutorial we'll find out! Join along in implementing an event-driven Natural Language Processing tool that does streaming HTML parsing, entity extraction, and sentiment analysis.

Just here for the code? Check it out [here](https://github.com/rotationalio/ensign-examples/tree/main/go/nlp)!

## Prerequisites

To follow along with this tutorial you'll need to:

- [Generate an API key to access Ensign]({{< ref "/getting-started#getting-started" >}})
- [Set up your GOPATH and workspace](https://go.dev/doc/gopath_code)
- [Create an Ensign client]({{< ref "/getting-started#create-a-client" >}})

## Back to the Future

Did you know? Some of the [earliest deployed machine learning apps](https://en.wikipedia.org/wiki/Naive_Bayes_spam_filtering) were event-driven! That’s right — back in the 90’s, email spam filters used Bayesian models to learn on the fly.

Spam filtering is an awesome example of a natural use case for [online modeling](https://en.wikipedia.org/wiki/Online_machine_learning). Each newly flagged spam message was a new training event, an opportunity to update the model in real time. While most machine learning bootcamps teach us to expect data in batches, there are a TON of natural use cases for streaming data science (maybe even more than for [offline aka batchwise modeling](https://en.wikipedia.org/wiki/Offline_learning)!).

Another great use case for event-driven data science is Natural Language Processing tasks such as named entity recognition, sentiment analysis, and text classification. In this tutorial, we'll tap into a live data feed and see how to process the text content as it streams in.

## A Whale of a Problem

The data we're going to be working with in this tutorial comes from a live RSS feed ingestion engine called [Baleen](https://github.com/rotationalio/baleen).

![baleen_diagram](/img/baleen_diagram.png)

If you want to run your own Baleen, check out [this branch](https://github.com/rotationalio/baleen/tree/ensign-demo). Then install the Baleen CLI:

```bash
$ go install ./cmd/baleen/
```

Then you can add posts with

```bash
$ baleen posts:add https://www.news-source-of-your-choice.com/link-to-article
```

Baleen has an Ensign `Publisher` that emits new events to a topic stream (let's call it `"baleen-docs"`) every time a new article is ingested. We should first check that we can access that particular Ensign stream (*Note: make sure you [create an Ensign `client`]({{< ref "/getting-started#create-a-client" >}}) first!*):

```golang
import (
    // ...
    ensign "github.com/rotationalio/go-ensign"
)

// This is the nickname of the topic, which gets mapped to an ID
// that actually gets used by Ensign to identify the stream
const Baleen = "baleen-docs"

func main() {
    // ...

	// Check to see if topic exists and create it if not
	exists, err := client.TopicExists(context.Background(), Baleen)
	if err != nil {
		panic(fmt.Errorf("unable to check topic existence: %s", err))
	}

	var topicID string
	if !exists {
		if topicID, err = client.CreateTopic(context.Background(), Baleen); err != nil {
			panic(fmt.Errorf("unable to create topic: %s", err))
		}
	} else {
		topics, err := client.ListTopics(context.Background())
		if err != nil {
			panic(fmt.Errorf("unable to retrieve project topics: %s", err))
		}

		for _, topic := range topics {
			if topic.Name == Baleen {
				var topicULID ulid.ULID
				if err = topicULID.UnmarshalBinary(topic.Id); err != nil {
					panic(fmt.Errorf("unable to retrieve requested topic: %s", err))
				}
				topicID = topicULID.String()
			}
		}
	}

    // ...
}
```

We can write a `Subscriber` to connect to the Baleen topic feed in order to tap into the `baleen-docs` stream:

```golang
import (
    // ...
    ensign "github.com/rotationalio/go-ensign"
)

func main() {
    // ...

    var sub ensign.Subscriber

	// Create a downstream consumer for the event stream
	sub, err = client.Subscribe(context.Background(), topicID)
	if err != nil {
		panic(fmt.Errorf("could not create subscriber: %s", err))
	}
	defer sub.Close()

    // ...
}
```

Next, you'll want to create a channel to consume events from the stream:

```golang
import (
    // ...
    api "github.com/rotationalio/go-ensign/api/v1beta1"
)

    // still inside main():

    var events <-chan *api.Event

    if events, err = sub.Subscribe(); err != nil {
        panic("failed to create subscribe stream: " + err.Error())
    }

    // ...
```

Now we've got the data feed ready, and the next step is create a loop that will listen on the `events` channel, and for each event it retrieves, do some NLP magic.

## NLP Magic Time

Ok, in the last section we set up our Ensign client and connected an Ensign `Subscriber` to Baleen's Ensign `Publisher`.

```golang
import (
    // ...
    post "github.com/rotationalio/baleen/events"
)

    // still inside main():

    // Events are processed as they show up on the channel
    for event := range events {
        if event.Type.Name == "Document" {

            // Unmarshal the event into an HTML Document
            doc := &post.Document{}
            if _, err = doc.UnmarshalMsg(event.Data); err != nil {
                panic("failed to unmarshal event: " + err.Error())
            }

            // Do NLP magic here!

        }
```

### Parsing the Beautiful Soup

The first step in all real world text processing and modeling projects (well, after ingestion of course ;-D) is parsing. The specific parsing technique has a lot to do with the data; but in this case we're starting with HTML documents, which is what Baleen's `Publisher` delivers.

Most data scientists are probably used to using [BeautifulSoup](https://www.crummy.com/software/BeautifulSoup/bs4/doc/) for HTML document parsing (just as long as it's [not regex!](https://stackoverflow.com/questions/1732348/regex-match-open-tags-except-xhtml-self-contained-tags/1732454#1732454)). For those doing their data science tasks in Golang, check out [Anas Khan](https://github.com/anaskhan96)'s [soup](https://github.com/anaskhan96/soup) package. Like the original Python package, it has a lot of great utilities for retrieving and preprocessing HTML.

Since we already have our HTML (which is in the form of bytes), all we need to do to prepare the soup and grab all the paragraphs using the `<p>` tags:

```golang
doc := soup.HTMLParse(string(html_bytes))
paras := doc.FindAll("p")
```

Now we can iterate over `paras` to process each paragraph chunk by chunk.


### More than a Feeling

Let's say that we want to do streaming sentiment analysis so that we can gauge the sentiment levels of the documents *right away* rather than in a batch analysis a month from now, when it may be too late to intervene!

For this we'll leverage the sentiment analysis tools implemented in [Connor DiPaolo](https://github.com/cdipaolo)'s [Golang Sentiment Library](https://github.com/cdipaolo/sentiment).

First we load the pre-trained model (trained using a Naive Bayes classifier, if you're curious!) using the `Restore` function:


```golang
import (
    "github.com/cdipaolo/sentiment"
)

    // ...

    var model sentiment.Models

    if model, err = sentiment.Restore(); err != nil {
        fmt.Println("unable to load pretrained model")
    }

    // ...
```

Now, let's iterate over the `paras` we extracted from the HTML in the section above and score the text of each using the pre-trained sentiment model:

```golang

    // ...

    var sentimentScores []uint8

    for _, p := range paras {

        // Get the sentiment score for each paragraph
        analysis := model.SentimentAnalysis(p.Text(), sentiment.English)
        sentimentScores = append(sentimentScores, analysis.Score)

        // ... more magic coming soon!
    }
```

We could look at the sentiment of each paragraph, but for tutorial purposes we'll just take an average sentiment for the overall article:

```golang
    // Get the average sentiment score across all the paragraphs
    var total float32 = 0
    for _, s := range sentimentScores {
        total += float32(s)
    }
    avgSentiment = total / float32(len(sentimentScores))
    }
```

But think of all the other things we can do with all that text!

### Finding the Who, Where, and What

Let's add an entity extraction step to our iteration over the `paras`. For this we need another dependency, the [prose](https://github.com/jdkato/prose) library created by [Joseph Kato](https://github.com/jdkato).

```golang
import (
    // ...
	prose "github.com/jdkato/prose/v2"
)

	// allocate empty entity map
	entities = make(map[string]string)

	for _, p := range paras {

		// Get sentiment score

		// Parse out the entities
		var parsed *prose.Document
		if parsed, err = prose.NewDocument(p.Text()); err != nil {
			fmt.Println("unable to parse text")
		}
		for _, ent := range parsed.Entities() {
			// add entities to map
			entities[ent.Text] = ent.Label
		}
	}
```

For those familiar with the Python library [spaCy](https://spacy.io/), `prose` works in a similar fashion. You first create a `prose.Document` by passing in the text content, which invokes the entity parsing. You can then iterate over the resulting `Entities`, which consist of tuples of the form `(Text, Label)`.

Take for example the sentence:
> Robyn Rihanna Fenty, born February 20, 1988, is a Barbadian singer, actress, and businesswoman.

The resulting entities and labels will be:

```
{
  "Barbadian": "GPE",
  "Robyn Rihanna Fenty": "PERSON"
}
```

Love you, Riri! Thanks to `soup`, `sentiment`, and `prose` (as well as Ensign and Baleen) we now have:

 - a live feed of RSS articles
 - a way to parse incoming HTML text into component parts
 - a way to score the sentiment of incoming articles
 - a way to extract entities from those articles

### What's Next?

So many possibilities! We could create a live alerting system that throws a flag every time a specific entity is mentioned. We could configure those alerts to fire only when the sentiment is below some threshold. Reach out to us at info@rotational.io and let us know what else you'd want to make!


## Breaking Free from the Batch

Applied machine learning has come a loooong way in the last ten years. Open source libraries like [scikit-learn](https://scikit-learn.org/stable/), [TensorFlow](https://www.tensorflow.org/), [spaCy](https://spacy.io/), and [HuggingFace](https://huggingface.co/) have put ML into the hands of everyday practitioners like us. However, many of us are still struggling to get our models into production.

And if you know how applied machine learning works, you know delays are bad! As new data naturally "drifts" away from historic data, the training input of our models becomes less and less relevent to the real world problems we're trying to use prediction to solve. Imagine how much more robust your applications would be if they were not only trained on the freshest data, but could alert you to drifts *as soon as they happen* -- you'd be able to react immediately as opposed to a batchwise process where you'd be lucky to catch the issue within a day!

Event-driven data science is one of the best solutions to the MLOps problem. MLOps often requires us to shoehorn our beautiful models into the existing data flows of our organizations. With a few very special exceptions (we especially love [Vowpal Wabbit](https://vowpalwabbit.org/) and [Chip Huyen's introduction to streaming for data scientists](https://huyenchip.com/2022/08/03/stream-processing-for-data-scientists.html)), ML tools and training teach us to expect our data in batches, but that's not usually how data flows organically through an app or into a database. If you can figure out how to reconfigure your data science flow to more closely match how data travels in your organization, the pain of MLOps can be reduced to almost nil.

Happy Eventing!

