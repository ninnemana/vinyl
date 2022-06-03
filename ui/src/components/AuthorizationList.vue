<template>
	<v-list>
		<v-list-item v-if='authorized'>
			<v-list-item-title>
				<router-link to='/account'>My Account</router-link>
			</v-list-item-title>
		</v-list-item>
		<v-list-item v-if='authorized'>
			<v-list-item-title>
				<router-link to='/logout'>Logout</router-link>
			</v-list-item-title>
		</v-list-item>
		<v-list-item v-if='!authorized'>
			<v-list-item-title>
				<router-link to='/login'>Login</router-link>
			</v-list-item-title>
		</v-list-item>
	</v-list>
</template>

<script lang='ts'>
import { Vue, Component, Prop } from 'vue-property-decorator';

@Component
export default class AuthorizationList extends Vue {
	@Prop(Boolean) readonly authorized!: boolean | false;
}
</script>

<!-- Add 'scoped' attribute to limit CSS to this component only -->
<style scoped lang='scss'>
</style>
