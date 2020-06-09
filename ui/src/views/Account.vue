<template>
    <v-content>
        <v-container fluid v-if='HasAccount'>
            <h1>Account Details</h1>
            <v-form>
                <v-container>
                    <v-row justify="center">
                        <v-col cols="12" md="8">
                            <v-text-field 
                                v-model="account.user.name" 
                                :counter="10" 
                                label="Name" 
                                required
                                center
                            />
                        </v-col>
                    </v-row>
                    <v-row justify="center">
                        <v-col cols="12" md="8">
                            <v-text-field 
                                v-model="account.user.email" 
                                :rules="emailRules()" 
                                label="Email address" 
                                required 
                            />
                        </v-col>
                    </v-row>
                </v-container>
            </v-form>
        </v-container>
    </v-content>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';

@Component({
	components: {}
})
export default class Account extends Vue {
    private account: any;

    constructor() {
        super();
        this.account = null;
    }

    emailRules(): any {
        return [ 
            (v: any) => !v || /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(v) || 'E-mail must be valid'
        ];
    }

    computed() {
        return {
            account: this.$store.getters.getAccount,
        };
    }

    get HasAccount(): boolean {
        const acct = this.$store.getters.getAccount;
        if (acct === null) {
            return false;
        }

        this.account = acct;
        return true;
    }
}
</script>
