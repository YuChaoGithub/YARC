import Vue from 'vue';
import VueRouter from 'vue-router';
import Home from './components/page/Home';
import Subreddit from './components/page/Subreddit';
import CreateArticle from './components/page/CreateArticle';
import ViewArticle from './components/page/ViewArticle';
import CreateSubreddit from './components/page/CreateSubreddit';
import UserProfile from './components/page/UserProfile';

Vue.use(VueRouter);

export const router = new VueRouter({
  mode: 'history',
  routes: [
    {
      name: 'home',
      path: '/',
      component: Home
    },
    {
      name: 'submit',
      path: '/submit',
      component: CreateArticle
    },
    {
      name: 'subreddit',
      path: '/r/:subreddit',
      component: Subreddit
    },
    {
      name: 'article',
      path: '/r/:subreddit/p/:pid',
      component: ViewArticle
    },
    {
      name: 'create_subreddit',
      path: '/subreddits/create',
      component: CreateSubreddit
    },
    {
      name: 'user',
      path: '/u/:username',
      component: UserProfile
    }
  ]
});
