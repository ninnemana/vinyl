import Vue from 'vue'
import VueRouter from 'vue-router'
import store from '../store'
import Account from '../views/Account.vue'
import CreateAccount from '../views/CreateAccount.vue'
import Home from '../views/Home.vue'
import Login from '../views/Login.vue'
import ResetPassword from '../views/ResetPassword.vue'

Vue.use(VueRouter)

const routes = [
	{
		path: '/',
		name: 'Home',
		component: Home
	},
	{
		path: '/login',
		name: 'Login',
		component: Login
	},
	// {
	// 	path: '/logout',
	// 	name: 'Logout',
	// 	beforeEnter: async () => {
	// 		await store.dispatch('logout');
	// 	},
	// 	redirect: { name: 'Login' }
	// },
	{
		path: '/account',
		name: 'Account',
		component: Account
	},
	{
		path: '/account/create',
		name: 'Create Account',
		component: CreateAccount
	},
	{
		path: '/account/reset-password',
		name: 'Reset Password',
		component: ResetPassword
	},
	{
		path: '/search',
		name: 'search',
		// route level code-splitting
		// this generates a separate chunk (about.[hash].js) for this route
		// which is lazy-loaded when the route is visited.
		component: () => import(/* webpackChunkName: "about" */ '../views/Search.vue'),
		beforeEnter: (to: any, from: any, next: any) => {
			store.dispatch('search', to.query.query);
			next();
		},
	},
	{
		path: '/release/:title/:id',
		name: 'details',
		// route level code-splitting
		// this generates a separate chunk (about.[hash].js) for this route
		// which is lazy-loaded when the route is visited.
		component: () => import(/* webpackChunkName: "about" */ '../views/Details.vue'),
		beforeEnter: (to: any, from: any, next: any) => {
			store.dispatch('get', to.params.id);
			next();
		},
	},
]

const router = new VueRouter({
	mode: 'history',
	base: process.env.BASE_URL,
	routes
})

export default router
