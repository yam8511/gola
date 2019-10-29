var vm = new Vue({
    el: '#app',
    data: {
        // 設置資訊
        showSetup: false,
        combineSetup: {},
        combine: [],
        playerCount: 0,
        targetPoint: 0,
        advanced: false,
        basicCard: {},
        selectedCards: {},
        selectedCombine: '',
        gameover: false,
        myName: '',

        // 桌面資訊
        PlayerPoint: {},
        PlayerName: {},
        ThrowCard: {},
        CardNum: {},
        MyCard: {},
        GameResult: '',

        needPlayCard: false,
        PointChanged: {},
        ShowResult: false,

        // 顯示出牌
        showPlayCard : false,
        PlayPlayerName : '',
        PlayCard : {},
        TurnMe: false,

        showOtherCard: false,
        showDrawedCard: false,
        showTakeCard: false,
        otherCard: [],
        GetCard: {},
        DrawedCard: {},
        needDrawCard: false,
        drawCard: [],
        hasDrawedCard: false,

        // 以下不知道會不會再用到
        list: '',
        player: '',
        showNumber: false,
        numbers: [],
        selectedNumber: 0,
        gameStart: false,
        gameOver: false,
        needEnterSession: false,
        isWaiting: false,
        shareToken: false,
        waitingOptions: {},
        onlyOneOption: false,
        isOut: false,
        online: [],
        intervalID: 0,
        ruleVisible: true,
        setupRules: {},

        // 系統用
        no: '',
        leader: false,
        ws: null,
        synth: null,
        voices: [],
        voice: null,
        token: '',
        session: '',
        input: '',
        action: '',
        output: '',
        is_mute: true,
        initializing: true,
        chatting: false,
        chatDone: false,
        chatOK: false,
        userEmail: '',
        suggest: '',
    },
    methods: {
        chat() {
            if (this.suggest.trim() === '') {
                return
            }
            this.chatDone = false
            this.chatting = true
            axios.post('/api/suggest', {
                email: this.userEmail,
                suggest: this.suggest,
                game: 'cd',
            })
                .then((response) => {
                    this.chatDone = true

                    if (response.status !== 200 || response.data !== 'ok') {
                        this.chatOK = false
                        return
                    }

                    this.userEmail = ''
                    this.suggest = ''
                    this.chatOK = true
                })
                .catch((error) => {
                    console.log('suggest error ->', error)
                }).finally(() => {
                    if (this.chatOK) {
                        this.chatting = false
                    } else {
                        this.chatting = true
                        setTimeout(() => {
                            this.chatting = false
                        }, 10000)
                    }
                })
        },
        joinGame() {

            if (this.ws !== null) {
                if (this.needEnterSession) {
                    this.input = JSON.stringify({
                        name: this.myName,
                        no: parseInt(this.session),
                    })
                    this.sendGame()
                    return
                }

                return
            }

            const ws = new WebSocket(this.wsURL + '?token=' + this.token)
            this.ws = ws
            this.intervalID = setInterval(() => {
                if (this.ws !== null) {
                    this.ws.send('ping')
                }
            }, 10000)
            ws.onmessage = this.receiveGame
            ws.onclose = this.closedGame
        },
        closedGame() {
            clearInterval(this.intervalID)
            this.ws = null
            this.resetInfo()
            if (this.gameStart) {
                this.output = '遊戲已經關閉'
            }
            this.gameStart = false
        },
        sendGame() {
            this.ws.send(JSON.stringify({
                'action': this.action,
                'reply': `${this.input}`,
            }))
        },
        changeSetup(e) {
            const name = e.target.id
            const n = parseInt(e.target.value)
            if (n === 0 || n === null) {
                this.selectedCards[name] = 0
            } else {
                this.selectedCards[name] = n
            }

            let extra = 0
            for (const name in this.selectedCards) {
                if (name === '隨機') {
                    continue
                }
                extra += this.selectedCards[name]
            }

            extra = this.playerCount * 4 - extra
            this.selectedCards['隨機'] = extra
            this.selectedCards = { ...this.selectedCards }

            this.input = JSON.stringify({
                PlayerCount: this.playerCount,
                TargetPoint: this.targetPoint,
                Advanced: this.advanced,
                CardSet: this.selectedCards,
            })
        },
        selectNumber(n) {
            if (this.selectedNumber !== 0) {
                return
            }

            this.selectedNumber = n
            this.numbers = [this.selectedNumber]
            this.input = JSON.stringify({
                'no': n,
                'name': this.myName,
            })
            this.sendGame()
        },
        selectOptions(k, v) {
            this.waitingOptions = {}
            if (!this.onlyOneOption) {
                this.waitingOptions[k] = v
            }
            this.waitingOptions = {...this.waitingOptions}
            this.input = v
            this.sendGame()
        },
        resetInfo() {
            this.leader = false
            this.no = ''
            this.session = ''
        },
        finishGame() {
            this.gameOver = false
            this.output = ''
        },
        changeLang(v) {
            this.voice = v.voice
            M.Toast.dismissAll();
            setTimeout(() => {
                M.toast({ html: `改成 ${v.voice.name} 發音` })
            }, 1)
        },
        lookCombineDetail(rules) {
            let display = []
            for (const key in rules) {
                display.push(`${key} x ${rules[key]}`)
            }

            return display.join(' | ')
        },
        changeMute() {
            this.is_mute = !this.is_mute
            this.speak(' ')
        },
        speak(sound) {
            if (sound === '' || sound === null || this.is_mute) {
                return
            }
            const speech = new window.SpeechSynthesisUtterance(sound)
            speech.voice = this.voice
            window.speechSynthesis.speak(speech)
        },
        playCard(index) {
            this.needPlayCard = false
            this.hasDrawedCard = true
            this.input = index
            this.sendGame()
        },
        receiveGame(e) {
            if (e.data === 'ping') {
                return
            }

            const tmp = JSON.parse(e.data)

            switch (tmp.action) {
                case 'refresh_online':
                    this.online = tmp.data
                    return

                case 'change_room_master':
                    this.leader = true
                    this.gameStart = tmp.data['遊戲開始']
                    if (!this.gameStart) {
                        this.action = 'change_room_master'
                        this.input = 'start'
                    }
                    return

                case 'refresh_desktop':
                    this.gameStart = true

                    this.playerCount = tmp.data.PlayerNo.length
                    this.targetPoint = tmp.data.TargetPoint

                    this.PlayerPoint = tmp.data.PlayerPoint
                    this.PlayerName = tmp.data.PlayerName
                    this.CardNum = tmp.data.CardNum

                    this.ThrowCard = {}
                    for (const no in tmp.data.ThrowCard) {
                        const cards = tmp.data.ThrowCard[no]
                        this.ThrowCard[no] = cards.map((card) => ({
                            ...card,
                            Detail: card.Detail.trim().split("\n\t").join('<br />'),
                        }))
                    }
                    this.ThrowCard = { ...this.ThrowCard }

                    this.MyCard = tmp.data.MyCard.map((card) => ({
                        ...card,
                        Detail: card.Detail.trim().split("\n\t").join('<br />'),
                    }))
                    return
            }




            if (!(tmp.action === 'game_is_running' && !isNaN(this.getHash()))) {
                if (tmp.sound !== '') {
                    this.output = tmp.sound
                }

                if (tmp.display !== '') {
                    this.output = tmp.display
                }

                this.speak(tmp.sound)
            }


            this.shareToken = false
            this.showSetup = false
            this.showNumber = false
            this.needEnterSession = false
            this.isWaiting = false
            this.waitingOptions = { '確認': '' }
            this.ShowResult = false
            this.showPlayCard = false
            this.showOtherCard = false
            this.showTakeCard = false
            this.showDrawedCard = false
            this.needDrawCard = false
            this.onlyOneOption = true
            this.action = tmp.action

            switch (tmp.action) {
                case 'card_setup':
                    // 快速組合
                    this.combineSetup = tmp.data.Combine
                    this.combine = Object.keys(this.combineSetup)
                    // 基本牌
                    this.basicCard = tmp.data.BasicCard
                    // 選取的卡牌
                    this.selectedCards = {}
                    this.selectedCombine = ''

                    const key = this.combine.length > 0 ? this.combine[0] : '';

                    if (key !== '') {
                        this.selectedCombine = key
                        const c = this.combineSetup[key]
                        for (const cardName in c.CardSet) {
                            this.selectedCards[cardName] = c.CardSet[cardName]
                        }
                        this.playerCount = c.PlayerCount
                        this.targetPoint = c.TargetPoint
                        this.advanced = c.Advanced

                        this.input = JSON.stringify({
                            PlayerCount: this.playerCount,
                            TargetPoint: this.targetPoint,
                            Advanced: this.advanced,
                            CardSet: this.selectedCards,
                        })
                    }

                    this.showSetup = true
                    this.gameOver = false
                    break

                case 'select_number':
                    this.resetInfo()
                    // 選擇號碼
                    this.selectedNumber = 0
                    this.showNumber = true
                    this.numbers = tmp.data
                    this.showNameInput = true
                    break

                case 'take_rule':
                    this.session = tmp.data['No']
                    this.no = tmp.data['No']
                    this.myName = tmp.data['Name']
                    this.output = '等待遊戲開始...'
                    window.location.hash = this.no
                    break

                case 'turn_me':
                    this.TurnMe = this.no === tmp.data
                    break

                case 'play_card':
                    this.MyCard = tmp.data.map((card) => ({
                        ...card,
                        Detail: card.Detail.trim().split("\n\t").join('<br />'),
                    }))
                    this.needPlayCard = true
                    break

                case 'show_play_card':
                    this.TurnMe = false
                    this.showPlayCard = true
                    this.PlayPlayerName = tmp.data.Player
                    this.PlayCard = {
                        ...tmp.data.Card,
                        Detail: tmp.data.Card.Detail.trim().split("\n\t").join('<br />'),
                    }
                    break

                case 'draw_card':
                    this.drawCard = tmp.data
                    this.needDrawCard = true
                    this.hasDrawedCard = false
                    break

                case 'look_player_cards':
                    this.showOtherCard = true
                    this.otherCard = tmp.data
                    this.isWaiting = true
                    break

                case 'show_result':
                    this.GameResult = tmp.data.GameResult
                    this.PlayerName = tmp.data.PlayerName
                    this.PlayerPoint = tmp.data.PlayerPoint
                    this.PointChanged = tmp.data.PointChanged
                    this.ShowResult = true
                    this.isWaiting = true
                    break

                case 'show_take_card':
                    this.GetCard = tmp.data.GetCard === null ? null : {
                        ...tmp.data.GetCard,
                        Detail: tmp.data.GetCard.Detail.trim().split("\n\t").join('<br />'),
                    }

                    this.MyCard = tmp.data.MyCard.map((card) => ({
                        ...card,
                        Detail: card.Detail.trim().split("\n\t").join('<br />'),
                    }))

                    this.isWaiting = true
                    this.showTakeCard = true

                    break

                case 'show_draw_card':
                    this.GetCard = tmp.data.GetCard === null ? null : {
                        ...tmp.data.GetCard,
                        Detail: tmp.data.GetCard.Detail.trim().split("\n\t").join('<br />'),
                    }
                    this.DrawedCard = tmp.data.DrawedCard === null ? null : {
                        ...tmp.data.DrawedCard,
                        Detail: tmp.data.DrawedCard.Detail.trim().split("\n\t").join('<br />'),
                    }
                    this.MyCard = tmp.data.MyCard.map((card) => ({
                        ...card,
                        Detail: card.Detail.trim().split("\n\t").join('<br />'),
                    }))

                    this.isWaiting = true
                    this.showTakeCard = true
                    this.showDrawedCard = true

                    break

                case 'select_player':
                    // 選擇號碼
                    this.selectedNumber = 0
                    this.showNumber = true
                    this.numbers = tmp.data
                    this.onlyOneOption = tmp.data.length == 1
                    break

                case 'game_is_running':
                    this.needEnterSession = true

                    let hash = this.getHash()
                    console.log('Hash ===> ', hash, tmp.data)
                    this.selectedNumber = 0
                    if (isNaN(hash) || tmp.data.indexOf(hash) < 0) {
                        this.showNumber = true
                        this.numbers = tmp.data
                        this.showNameInput = true
                    } else {
                        this.showNumber = true
                        this.numbers = [hash]
                        this.showNameInput = true
                        this.selectNumber(hash)
                    }

                    break
                case 'all_close_eyes':
                    this.gameStart = true
                    break
                case 'waiting':
                    this.isWaiting = true
                    if (tmp.data !== null && Object.keys(tmp.data).length > 0) {
                        this.onlyOneOption = Object.keys(tmp.data).length == 1
                        this.waitingOptions = {}
                        for (const key in tmp.data) {
                            this.waitingOptions[key] = tmp.data[key]
                        }
                    }
                    break
                case 'take_token':
                    this.token = tmp.data
                    window.location.search = this.token
                    // this.shareToken = true
                    // this.isWaiting = true
                    break
                case 'game_over':
                    this.gameStart = false
                    this.gameOver = true
                    this.isOut = false
                    this.resetInfo()
                    break
                default:
                    break
            }
        },
        getHash() {
            const hashes = window.location.hash.split('#', 2)
            let hash = ''
            if (hashes.length == 1) {
                hash = hashes[0]
            } else if (hashes.length > 1) {
                hash = hashes[1]
            }

            hash = parseInt(hash)
            return hash
        },
    },
    watch: {
        selectedCombine(v) {
            const setup = this.combineSetup[v]
            this.playerCount = setup.PlayerCount
            this.targetPoint = setup.TargetPoint
            this.advanced = setup.Advanced

            for (const name in this.basicCard) {
                this.selectedCards[name] = 0
            }

            for (const name in setup.CardSet) {
                const count = setup.CardSet[name]
                this.selectedCards[name] = count
            }

            this.input = JSON.stringify({
                PlayerCount: this.playerCount,
                TargetPoint: this.targetPoint,
                Advanced: this.advanced,
                CardSet: this.selectedCards,
            })
        },
    },
    computed: {
        qrcode() {
            return window.location.origin + window.location.pathname + '?' + this.token
        },
        homeQRCode() {
            return window.location.origin + "/cd"
        },
        gameRules() {
            return this.lookCombineDetail(this.setupRules)
        },
        optionsCol() {
            const optionLen = parseInt(12 / Object.keys(this.waitingOptions).length)
            return optionLen > 12 ? 1 : optionLen
        },
        wsURL() {
            let url = window.location.host + '/cd/game'
            if (window.location.protocol === 'http:') {
                url = 'ws://' + url
            } else {
                url = 'wss://' + url
            }
            return url
        }
    },
    mounted() {
        console.log('Ws URL', this.wsURL)
        const modelElems = document.querySelectorAll('.modal')
        let modelInstances = M.Modal.init(modelElems, {
            onCloseEnd: () => {
                this.chatDone = false
            }
        })

        if ('speechSynthesis' in window) {
            this.synth = window.speechSynthesis
            setTimeout(() => {
                let voices = this.synth.getVoices()
                const allowVoices = {
                    'zh-TW': true,
                    'zh-CN': true,
                    'zh-HK': true,
                    'zh-HK': true,
                    'en-US': true,
                    'ja-JP': true,
                    'ko-KR': true,
                }
                voices = voices.
                    filter((e) => allowVoices[e.lang] ? allowVoices[e.lang] : false).
                    map((e) => {
                        let words = e.name.split(' ')
                        let text = words[words.length - 1]
                        if (words.length == 1) {
                            words = e.name.split(' ')
                            text = words[words.length - 1]
                            words = text.split('（')
                            text = words[0]
                        }
                        if (e.lang === 'zh-TW') {
                            this.voice = e
                        }
                        return {
                            text: text,
                            lang: e.lang,
                            voice: e,
                        }
                    })

                this.voices = {}
                for (const i in voices) {
                    const voice = voices[i]
                    this.voices[voice.lang] = voice
                }

                setTimeout(() => {
                    const elems = document.querySelectorAll('.fixed-action-btn');
                    let instance = M.FloatingActionButton.init(elems, {
                        hoverEnabled: false,
                    });

                    this.initializing = false
                }, 500)
            }, 500)

            if (window.location.search !== '') {
                let hashes = window.location.search.split('?', 2)
                if (hashes.length == 1) {
                    this.token = hashes[0]
                } else if (hashes.length > 1) {
                    this.token = hashes[1]
                }

                this.joinGame()
            }
        }
    }
})
