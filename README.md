# unquomment

a absolute unnecessary api that returns stupid comments for a specified topic

The comments are generated by OpenAI GPT-3 which is promted to create a comment by a bored ~40 year old man in a meeting


## Usage

    curl --request GET \
      --url 'https://unquomment.hmnd.de/comment?topic=ai&sexism=3&boredom=1&stupidity=10&aggression=10'

## Parameters

- `topic` (required) - The topic of the comment
- `sexism` (optional) - The sexism of the comment (0-10)
- `boredom` (optional) - The boredom of the comment (0-10)
- `stupidity` (optional) - The stupidity of the comment (0-10)
- `aggression` (optional) - The aggression of the comment (0-10)