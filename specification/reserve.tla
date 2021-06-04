------------------------------- MODULE reserve -------------------------------
EXTENDS     Naturals, Sequences, Reals

CONSTANT    Coin,   \* Set of all coins
            Pair,   \* Set of all pairs of coins
            (* User constant is used to limit number of account records.   *)
            User,   \* Set of all users
           
VARIABLE    book,   \* Order Book
            bonds,  \* AMM Bond Curves
            
-----------------------------------------------------------------------------
NoVal ==    CHOOSE v : v \notin Nat

Amount == r \in Real

Inflation == r \in Real

(***************************************************************************)
(* The NOM coin is the representation of credit or a right to mint         *)
(* by the Onomy Reserve                                                    *)
(*                                                                         *)
(* Each user account has a single balance of NOM with potential for many   *)
(* outstanding balances of minted Denoms                                   *)
(*                                                                         *)
(* NOM: Credit                                                             *)
(* Denoms: Debits                                                          *)
(***************************************************************************)
Account == [
    nom: Amount, 
    denoms: {[denom: Coin, amount: Amount]}
]

Reserve == [
    nom: Amount
]

(***************************************************************************)
(* Denom Specific Parameters voted by NOM holders                          *)
(*                                                                         *)
(* catio: minimum minting collateralization ratio (denom specific)         *)
(* latio: liquidation collateralization ratio (denom specific)             *)
(* destatio: percentage of denom staked at validator (denom specific)      *)
(***************************************************************************)
DeParam == [denom: Coin, catio: Real, destatio: Real, flatio: Real]

(***************************************************************************)
(* NOM coupons are used as a tradable index of currencies held in a reserve*)
(* account.                                                                *)
(*                                                                         *)
(* Coupons are tied to specific accounts but are not permissioned          *)
(*                                                                         *)
(* The coupon is denominated in NOM and is redeemable for NOM when         *)
(* burned along with the proportional amount of indexed currencies.        *)
(*                                                                         *)
(* The goal of this feature is to allow for monetization of reserve        *) 
(* rewards without liquidating NOM collateral. It also allows others than  *)
(* the account holder to swap the basket index of currencies for NOM       *)
(* given they surrender the amount of swaps corresponding to the amount of *)
(* NOM redeemed                                                            *)
(*                                                                         *)
(* A coupon is the "coupon" of stripped bond that acts as a NOM put        *)
(* against  the basket of denoms minted by a reserve account with an       *)
(* inflationary coupon rate determined by the global percentage of NOM     *) 
(* staked in reserve control variable.                                     *)
(***************************************************************************)
Coupon == [user: User, amount: Real, denoms: {[denom: Coin, amount: Amount]}]

Type == /\  bonds \in [Pair -> [Coin -> Amount]]
        /\  coupons \in [User -> Token]
            (***************************************************************)
            (* Time is abstracted to a counter that increments during a    *) 
            (* “time” step. All other steps are time stuttering            *)
            (*                                                             *)
            (* In blockchain this corresponds to the block                 *)
            (*                                                             *)
            (* In asynchronous DAG, like with Equity protocol,             *)
            (* recurring processes relying on time, such as inflation,     *)
            (* will be triggered by a timer ran on correct nodes.          *)
            (*                                                             *)
            (* Timestamps of recurring DAG process will be the average of  *)
            (* reported times by nodes for each recurring process          *)
            (***************************************************************)
        /\  time \in Real
        /\  accounts \in [User -> Account]
        /\  params \in [Coin -> Param]

(***************************************************************************)
(* Credit NOM to User Account. Add r to balance.                           *)
(***************************************************************************)
Deposit(user) ==  /\ \E r \in Reals :
                        /\ 'accounts = [accounts EXCEPT ![user].nom = @ + r]
                        /\ UNCHANGED << bonds, tokens, time, params >>
(***************************************************************************)
(* Debit NOM from User Account. Minus r from balance.                      *)
(***************************************************************************)
Withdraw(user) == /\ \E r \in Reals : r < account[user].nom :
                        /\ ‘accounts = [accounts EXCEPT ![user].nom = @ - r]
                        /\ UNCHANGED << bonds, tokens, time, params >>

(* Mint denom and bond NOM *)
Mint(user) ==  
    LET nomBal == accounts[user].nom
    IN 
        (*********************** Qualifying Condition **********************)
        (* Nom balance of user account greater than 0                      *)
        (*******************************************************************)
        /\ nomBal > 0
        (*******************************************************************)
        (* There exists r in Reals such that : r is less than the user     *)
        (* account nom balance                                             *)
        (*******************************************************************)
        /\ \E r \in Reals : r < nomBal :
            /\ accounts' = [accounts EXCEPT ![user].nom = @ - r]
            /\ reserve' = reserve + 1
            /\ 
    
(***************************************************************************)
(* Burn denom and unbond NOM                                               *)
(* Burning Denoms is like a past time, it’s fun.  Users really like doing  *)
(* it because it allows them to unlock their NOM, which they want to stake *)
(* with validators rather than mint Denoms.                                *)                         
(***************************************************************************)
Burn(user) ==   (* There exists r in Reals such that : r is less than      *)
                (* amount of coupons  *)
                (*                                                         *)
                (***********************************************************)
               /\ coupons[user] != {}
                /\ CHOOSE coupon \in coupons[user] : { coupon. :coupon /in coupons}



