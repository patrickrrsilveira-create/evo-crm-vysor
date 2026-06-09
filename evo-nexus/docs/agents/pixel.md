# Pixel — Social Media

**Command:** `/pixel-social-media` | **Color:** yellow | **Model:** Sonnet

Pixel is the social media agent responsible for content creation, editorial calendar management, and cross-platform analytics. It covers YouTube, Instagram, LinkedIn, Twitter/X, Threads, and Bluesky — writing posts, threads, and carousels, analyzing engagement metrics, tracking audience growth, and generating unified performance reports across all platforms.

## When to Use

- You need to write a LinkedIn post, Twitter thread, or carousel
- You want to create or update a content calendar for the week/month
- You need a cross-platform analytics report or platform-specific report (YouTube, Instagram, LinkedIn)
- You want to repurpose a blog post or video into social media content
- You need recommendations on what's working and how to improve engagement

## Preloaded Skills

| Skill | What it does |
|-------|-------------|
| `social-post-writer` | Write posts for LinkedIn, Twitter/X, Threads, or Bluesky |
| `social-thread-writer` | Create multi-part threads for Twitter/X or LinkedIn |
| `social-carousel-writer` | Build LinkedIn carousel slide content |
| `social-hook-writer` | Craft attention-grabbing opening lines |
| `social-content-calendar` | Plan and organize a posting schedule |
| `social-content-strategy` | Define topic clusters and content mix |
| `social-content-repurposer` | Adapt one piece of content across multiple platforms |
| `social-content-pattern-analyzer` | Identify what content types and topics perform best |
| `social-performance-analyzer` | Analyze post-level metrics and engagement |
| `social-optimization-advisor` | Get concrete improvement recommendations |
| `social-platform-strategy` | Platform-specific tactical guidance |
| `social-audience-growth-tracker` | Track and analyze follower growth trends |
| `social-context` | Set up or update profile, voice, and audience context |
| `social-analytics-report` | Unified cross-platform analytics dashboard |
| `social-instagram-report` | Instagram-specific analytics report |
| `social-linkedin-report` | LinkedIn-specific analytics report |
| `social-youtube-report` | YouTube-specific analytics report |
| `int-instagram` | Query Instagram Graph API for metrics and insights |
| `int-linkedin` | Query LinkedIn API for profile and org stats |
| `int-youtube` | Query YouTube Data API for channel and video stats |

## Example Interactions

```
/pixel-social-media write a LinkedIn post about our new AI agents feature
/pixel-social-media create a content calendar for next week
/pixel-social-media youtube report — how did our channel perform this month?
```

## Routines

| Time | Routine | Command |
|------|---------|---------|
| Daily 18:00 | Social Analytics | `make social` |
| Friday 08:15 | Social Analytics Weekly | `make social` |
| 1st of month | Social Analytics Monthly | `make social` |

## Memory

Persistent memory at `.claude/agent-memory/pixel/`. Stores project context, user preferences, and feedback across sessions.
