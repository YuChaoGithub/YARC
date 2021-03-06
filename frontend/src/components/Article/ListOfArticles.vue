<template>
  <div class="col">
    <!-- Sort, Create Post -->
    <div class="q-pt-md q-pl-md q-pr-md text-white">
      <q-list dark bordered class="bg-grey-10 rounded-borders">
        <q-item class="q-mt-sm q-mb-sm">
          <q-toggle
            v-if="showSubscribedToggle && $store.state.auth"
            class="q-mr-md"
            color="blue"
            dark v-model="subscribed"
            val="xs"
            label="Subscribed Only"
          />
          <div class="q-gutter-sm">
            <q-radio dark v-model="sortBy" val="hot" label="Hot" />
            <q-radio dark v-model="sortBy" val="new" label="New" />
            <q-radio dark v-model="sortBy" val="old" label="Old" />
          </div>

          <template v-if="canCreatePost">
            <q-input @click="createPost('text')" class="col q-ml-lg q-mr-sm" dark clearable outlined dense standout label="Create Post" />
            <q-btn @click="createPost('image')" flat round color="grey" icon="insert_photo">
              <q-tooltip>Create Image Post</q-tooltip>
            </q-btn>
            <q-btn @click="createPost('link')" flat round color="grey" icon="insert_link">
              <q-tooltip>Create Link Post</q-tooltip>
            </q-btn>
          </template>
        </q-item>
      </q-list>
    </div>

    <!-- Article Lists -->
    <div class="q-pa-md text-white">
      <q-list dark bordered separator class="bg-grey-10 rounded-borders">
        <q-infinite-scroll @load="loadMoreArticles" :offset="250" ref="infiniteScroll">
          <article-entry
            v-for="article in articles"
            :key="article.articleID"
            :article="article"
            @click.native="articleClicked(article)"
          />
        </q-infinite-scroll>
        
        <!-- Error Indicator -->
        <q-banner v-if="errOccurred" class="text-white bg-red">
            Error fetching the article list, please try again.
        </q-banner>
        
        <!-- Empty Indicator -->
        <q-item v-if="!loading && articles.length == 0 && !errOccurred" class="column">
          <div class="q-ma-lg">
            <div class="text-h3 q-mb-md">Wow, such empty!</div>
            <div class="text-h6">{{emptyText}}</div>
          </div>
        </q-item>

      </q-list>
    </div>
  </div>
</template>

<script>
import ArticleEntry from '../article/ArticleEntry';
import ArticleService from '../../services/article';

export default {
  props: {
    emptyText: {
      type: String,
      default: ''
    },
    criterion: {
      type: String,
      default: ''
    },
    criterionKey: {
      type: String,
      default: ''
    },
    showSubscribedToggle: { // null, true, or false. (The "subscribed" button won't show.)
      type: Boolean,
      default: false
    },
    subreddit: {
      type: String,
      default: ''
    },
    canCreatePost: {
      type: Boolean,
      default: true
    },
    loadOnMount: {
      type: Boolean,
      default: true
    }
  },
  data() {
    return {
      articles: [],
      sortBy: "hot",
      articlesPerRequest: 10,
      errOccurred: false,
      loading: true,
      subscribed: false
    };
  },
  watch: {
    sortBy() {
      this.reloadArticles();
    },
    criterion() {
      this.reloadArticles();
    },
    criterionKey() {
      this.reloadArticles();
    },
    subscribed() {
      this.reloadArticles();
    }
  },
  mounted() {
      this.reloadArticles();
  },
  methods: {
    articleClicked(article) {
      this.$router.push({
        path: '/r/' + article.subreddit + '/p/' + article.articleID
      });
    },
    createPost(type) {
      this.$router.push({
        name: 'submit',
        query: {
          subreddit: this.subreddit,
          postType: type
        }
      });
    },
    async fetchArticleLists(after) {
      let auth = this.$store.state.auth;
      let fetchedArticles = await ArticleService.getList(this.sortBy, after, this.articlesPerRequest, this.criterion, this.criterionKey, this.subscribed, auth ? auth.authHeader : null);
      
      // Check if it is the end of the article list. If it is, stop the scroll.
      if (fetchedArticles.length < this.articlesPerRequest) {
        this.$refs.infiniteScroll.stop();
      }

      return fetchedArticles;
    },
    async reloadArticles() {
      // Refetch the articles since the sort criterion has changed.
      this.loading = true;
      this.articles = [];
      try {
        this.articles = await this.fetchArticleLists("");
        this.$refs.infiniteScroll.resume(); // Make the infinite scroller work again.
        this.errOccurred = false;
      } catch {
        this.errOccurred = true;
      }
      this.loading = false;
    },
    loadMoreArticles(_, done) {
      // Wait for the initial articles to be loaded.
      if (this.loading) {
        done();
        return;
      }

      // Fetch more articles.
      let after = this.articles.length == 0 ? "" : this.articles[this.articles.length-1].articleID;
      this.fetchArticleLists(after).then(fetchedArticles => {
          for (const art of fetchedArticles) {
            this.articles.push(art);
          }
          done();
      }).catch(() => {
        this.errOccurred = true;
        done(true);
      });
    }
  },
  components: {
    articleEntry: ArticleEntry
  }
}
</script>