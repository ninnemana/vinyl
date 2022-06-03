import axios, { AxiosRequestConfig, AxiosResponse } from 'axios'
// import * as Cookies from 'js-cookie'
import Vue from 'vue'
import Vuex from 'vuex'
import VuexPersist from 'vuex-persist'

Vue.use(Vuex)

const vuexLocalStorage = new VuexPersist({
	key: 'vinylshare', // The key to store the state on in the storage provider.
	storage: window.localStorage, // or window.sessionStorage or localForage
	// Function that passes the state and returns the state with only the objects you want to store.
	// reducer: state => state,
	// Function that passes a mutation and lets you decide if it should update the state in localStorage.
	// filter: mutation => (true)
})

export default new Vuex.Store({
	plugins: [vuexLocalStorage.plugin],
	state: {
		searchTerm: '',
		searchResults: [] as any[],
		searchPagination: {},
		getResult: {},
		searchError: {},
		getError: {},
		authorized: null,
		account: {} as any,
		loginError: '',
	},
	mutations: {
		SET_SEARCH_TERM(state: any, term: string) {
			state.searchTerm = term;
		},
		SET_SEARCH_RESULTS(state: any, res: any) {
			const results: any[] = [];
			res.forEach((result: any, i: number) => {
				if (i === 0) {
					state.searchPagination = result.pagination;
				}
				
				results.push(result.release);
			});
			state.searchResults = results;
		},
		SET_SEARCH_ERROR(state: any, err: Error) {
			state.searchError = err;
		},
		SET_SEARCH_PAGE(state: any, page: number) {
			state.searchPagination = page;
		},
		SET_GET_RESULT(state: any, result: any) {
			state.getResult = result;
		},
		SET_GET_ERROR(state: any, error: Error) {
			state.getError = error;
		},
		SET_ACCOUNT(state: any, account: any) {
			state.account = account;
		},
		SET_LOGIN_ERROR(state: any, error: Error) {
			state.loginError = error;
		},
	},
	actions: {
		setSearchPage({ commit }, page: number) {
			commit('SET_SEARCH_PAGE', page);
		},
		createAccount({ commit }, account: AccountRequest) {
			axios({
				method: 'post',
				url: '${process.env.VUE_APP_API_DOMAIN}/account/create',
				data: account,
			}).then((r: AxiosResponse) => {
				commit('SET_ACCOUNT', r.data);
			}).catch((error: { response: { data: {} } }) => {
				if (!error.response) {
					commit('SET_CREATE_ACCOUNT_ERROR', 'Failed to create account.');
					return;
				}

				commit('SET_CREATE_ACCOUNT_ERROR', error.response.data);
			});
		},
		search({ commit }, query: string) {
			commit('SET_SEARCH_TERM', query)
			
			axios({
				method: 'post',
				url: `${process.env.VUE_APP_API_DOMAIN}/vinyls/search`,
				data: { artist: query },
				headers: { 
					Authorization: `Bearer ${this.state.account.token}`,
				},
			} as AxiosRequestConfig).then((r: AxiosResponse) => {
				commit('SET_SEARCH_RESULTS', r.data.results);
			}).catch((error: { response: { data: {} } }) => {
				if (!error.response) {
					commit('SET_SEARCH_ERROR', 'Failed to retrieve results.');
					return;
				}

				commit('SET_SEARCH_ERROR', error.response.data);
			});
		},
		get({ commit }, id: string) {
			axios({
				method: 'get',
				url: `${process.env.VUE_APP_API_DOMAIN}/vinyls/${id}`,
				headers: {
					Authorization: `Bearer ${this.state.account.token}`,
				},
			}).then((r: AxiosResponse) => {
				commit('SET_GET_RESULT', r.data)
			}).catch((error: { response: { data: {} } }) => {
				if (!error.response) {
					commit('SET_GET_ERROR', 'Failed to retrieve result.');
					return;
				}

				commit('SET_GET_ERROR', error.response.data)
			});
		},
		login({ commit }, auth: { email: string; password: string }) {
			axios({
				method: 'post',
				url: `${process.env.VUE_APP_API_DOMAIN}/auth`,
				data: auth,
			}).then((r: AxiosResponse) => {
				commit('SET_ACCOUNT', r.data);
			}).catch((error: { response: { data: {} } }) => {
				if (!error.response || !error.response.data) {
					commit('SET_LOGIN_ERROR', 'Failed to login.');
					return;
				}

				commit('SET_LOGIN_ERROR', error.response.data);
			});
		},
		logout({ commit }) {
			commit('SET_ACCOUNT', null);
		},
	},
	modules: {
	},
	getters: {
		getResult: (state: any) => {
			return state.getResult;
		},
		searchResults: (state: any) => {
			return state.searchResults;
		},
		getAccount: (state: any) => {
			return state.account;
		},
		getLoginError: (state: any) => {
			return state.loginError;
		},
	}
})

export class AccountRequest {
	private username: string;
	private email: string;
	private password: string;
	private errors: Array<string>;

	constructor(username: string, email: string, password: string) {
		this.username = username;
		this.email = email;
		this.password = password;
		this.errors = [];
	}

	valid(): boolean {
		this.errors = [];
		
		if (this.username === '') {
			this.errors.push('Invalid username');
		}

		if (this.email == '') {
			this.errors.push('Invalid e-mail address');
		}

		if (this.password == '') {
			this.errors.push('Invalid password');
		}

		return this.errors.length === 0;
	}
}
