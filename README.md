# Gator - RSS Feed Aggregator

A command-line RSS feed aggregator that helps you follow and aggregate RSS feeds.

## Setup

1. Install Go 1.24.0 or later
2. Clone the repository
3. Create a PostgreSQL database
4. Create a `.gatorconfig.json` file in your home directory:

```json
{
  "db_url": "postgresql://username:password@localhost:5432/gator",
  "current_username": ""
}
```

## Available Commands

### User Management

- `register <username>` - Create a new user account
- `login <username>` - Login as an existing user
- `users` - List all registered users
- `reset` - Delete all users from database

### Feed Management

- `addfeed <name> <url>` - Add a new RSS feed and follow it
- `feeds` - Show all available feeds
- `follow <feed_url>` - Follow an existing feed
- `unfollow <feed_url>` - Unfollow a feed
- `following` - List feeds you're following
- `browse [options]` - Browse posts with pagination and sorting options
- `agg <duration>` - Start feed aggregation with specified interval

## Example Usage

### 1. Initial Setup

```bash
# Create config file with your database credentials
echo '{
  "db_url": "postgresql://username:password@localhost:5432/gator",
  "current_username": ""
}' > ~/.gatorconfig.json
```

### 2. User Management

```bash
# Register new account
gator register john

# List all users
gator users

# Login as different user
gator login alice
```

### 3. Feed Management

```bash
# Add and follow a new feed
gator addfeed "Tech News" "https://example.com/feed.xml"

# List all feeds
gator feeds

# View feeds you're following
gator following
```

### 4. Content Aggregation

```bash
# Start aggregator (runs every minute)
gator agg 1m

# Browse posts with pagination and sorting
gator browse --page 1 --limit 5 --sort-by title --sort-order asc

# Available browse options:
# --page <number>      : Page number (default: 1)
# --limit <number>     : Posts per page (default: 2)
# --sort-by <field>    : Sort by field (title|url|published_at) (default: published_at)
# --sort-order <order> : Sort order (asc|desc) (default: desc)
```

### 5. Advanced Operations

```bash
# Unfollow a feed
gator unfollow "https://example.com/feed.xml"

# Delete all users (caution!)
gator reset
```
