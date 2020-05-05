import Vue from 'vue'
import Vuex from 'vuex'
import axios, { AxiosResponse } from 'axios'

Vue.use(Vuex)

export default new Vuex.Store({
	state: {
		searchTerm: '',
		searchResults: [] as any[],
		searchPagination: {},
		getResult: {},
		searchError: {},
		getError: {},
		authorized: null,
		account: {},
	},
	mutations: {
		SET_SEARCH_TERM: (state, term) => {
			state.searchTerm = term;
		},
		SET_SEARCH_RESULTS: (state, res) => {
			const results: any[] = [];
			res.forEach((result: any, i: number) => {
				if (i === 0) {
					state.searchPagination = result.pagination;
				}
				
				results.push(result.release);
			});
			state.searchResults = results;
		},
		SET_SEARCH_ERROR: (state, error) => {
			state.searchError = error;
		},
		SET_SEARCH_PAGE: (state, page) => {
			state.searchPagination = page;
		},
		SET_GET_RESULT: (state, result) => {
			state.getResult = result;
		},
		SET_GET_ERROR: (state, error) => {
			state.getError = error;
		},
		SET_ACCOUNT: (state, account) => {
			state.account = account;
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
				url: `${process.env.VUE_APP_API_DOMAIN}/search`,
				data: { artist: query },
			}).then((r: AxiosResponse) => {
				commit('SET_SEARCH_RESULTS', r.data);
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
			}).then((r: AxiosResponse) => {
				commit('SET_GET_RESULT', r.data);
			}).catch((error: { response: { data: {} } }) => {
				if (!error.response) {
					commit('SET_GET_ERROR', 'Failed to retrieve result.');
					return;
				}

				commit('SET_GET_ERROR', error.response.data);
			});
		},
	},
	modules: {
	},
	getters: {
		getResult: (state) => {
			return state.getResult;
		},
		searchResults: (state) => {
			return state.searchResults;
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
