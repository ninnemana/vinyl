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
	},
	actions: {
		setSearchPage({ commit }, page: number) {
			commit('SET_SEARCH_PAGE', page);
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
