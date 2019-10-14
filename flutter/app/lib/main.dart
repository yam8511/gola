import 'dart:async';

import 'package:flutter/material.dart';
import 'package:web_socket_channel/io.dart';
import 'package:app/helper/tts_helper.dart';
import 'package:flare_flutter/flare_actor.dart';

void main() => runApp(MyApp());

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      initialRoute: "/",
      title: '狼人殺',
      routes: {
        '/': (ctx) => LaunchPage(),
        'home': (ctx) => HomePage(),
        'setting': (ctx) => SettingPage(),
        'game': (ctx) => GamePage(),
      },
    );
  }
}

class LaunchPage extends StatefulWidget {
  @override
  _LaunchPageState createState() => _LaunchPageState();
}

class _LaunchPageState extends State<LaunchPage> {
  goPage() {
    Navigator.pushReplacementNamed(context, "home");
  }

  autoGoPage() async {
    var _duration = Duration(milliseconds: 2500);
    return new Timer(_duration, goPage);
  }

  @override
  void initState() {
    super.initState();
    autoGoPage();
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Color(0xFFFFFFFF),
      child: Column(
        children: <Widget>[
          Padding(
            padding: EdgeInsets.all(100.0),
          ),
          Container(
            color: Color(0xFFFFFFFF),
            height: 300,
            child: FlareActor(
              "assets/wolf.flr",
              alignment: Alignment.center,
              fit: BoxFit.contain,
              animation: "wolf",
            ),
          ),
          Container(
            alignment: Alignment.topCenter,
            color: Color(0xFFFFFFFF),
            height: 300,
            child: FlareActor(
              "assets/wolf_pen.flr",
              alignment: Alignment.center,
              fit: BoxFit.contain,
              animation: "Untitled",
            ),
          ),
        ],
      ),
    );
  }
}

class HomePage extends StatefulWidget {
  @override
  _HomePageState createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  TextEditingController _room_input = new TextEditingController();
  // IOWebSocketChannel _websocket_channel;
  int _sliderValue = 0;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          "狼人殺",
          style: TextStyle(color: Colors.white),
        ),
        backgroundColor: Colors.black87,
      ),
      body: Container(
        child: GestureDetector(
          behavior: HitTestBehavior.translucent,
          onTap: () {
            // 手指觸碰頁面收鍵盤
            FocusScope.of(context).requestFocus(FocusNode());
          },
          child: Column(
            children: <Widget>[
              Padding(
                padding: EdgeInsets.all(50.0),
              ),
              RaisedButton(
                padding: EdgeInsets.only(
                    left: 100.0, right: 100.0, top: 10.0, bottom: 10.0),
                color: Colors.orange,
                textColor: Colors.white,
                child: Text("建立房間"),
                shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(5.0)),
                onPressed: () {
                  Navigator.pushNamed(context, "setting");
                },
              ),
              Padding(
                padding: EdgeInsets.all(50.0),
              ),
              Divider(),
              Padding(
                padding: EdgeInsets.symmetric(
                  vertical: 16.0,
                  horizontal: 32.0,
                ),
                child: TextField(
                  controller: _room_input,
                  decoration: InputDecoration(
                    focusedBorder: UnderlineInputBorder(
                      borderSide: BorderSide(
                        color: Colors.orange,
                      ),
                    ),
                    labelText: "房間號碼",
                    hintText: "number",
                    labelStyle: TextStyle(color: Colors.orange),
                    icon: Icon(
                      Icons.home,
                      color: Colors.orange,
                    ),
                  ),
                ),
              ),
              RaisedButton(
                padding: EdgeInsets.only(
                    left: 100.0, right: 100.0, top: 10.0, bottom: 10.0),
                color: Colors.orange,
                textColor: Colors.white,
                child: Text("加入房間"),
                shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(5.0)),
                onPressed: () {
                  int room = int.parse(_room_input.text);
                  if (room != 0) {
                    Navigator.pushNamed(context, "game");
                  }
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class SettingPage extends StatefulWidget {
  @override
  _SettingPageState createState() => _SettingPageState();
}

class _SettingPageState extends State<SettingPage> {
  TextEditingController _token_input = new TextEditingController();

  var mapPlayerSelect = {
    "Human" : 0,
    "Wolf" : 0,
    "Knight" : 0,
    "Prophesier" : 0,
    "Hunter" : 0,
    "Witch" : 0,
    };
  
  int _sliderHuman = 0;
  int _sliderWolf = 0;
  int _sliderKnight = 0;
  int _sliderProphesier = 0;
  int _sliderHunter = 0;
  int _sliderWitch = 0;

  IOWebSocketChannel _websocket_channel;

  String host = "127.0.0.1";
  String gamestart = "/wf/";
  String gamename = "/wf/game";
  String ws = "ws://";
  String wss = "wss://";

  String _selectNum = "";

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          "設定",
          style: TextStyle(color: Colors.white),
        ),
        backgroundColor: Colors.black87,
      ),
      body: Container(
        child: GestureDetector(
          behavior: HitTestBehavior.translucent,
          onTap: () {
            FocusScope.of(context).requestFocus(FocusNode());
          },
          child: Column(
            children: <Widget>[
              Padding(
                padding: EdgeInsets.all(10.0),
              ),
              Padding(
                padding: EdgeInsets.symmetric(
                  vertical: 16.0,
                  horizontal: 32.0,
                ),
                child: TextField(
                  controller: _token_input,
                  decoration: InputDecoration(
                    focusedBorder: UnderlineInputBorder(
                      borderSide: BorderSide(
                        color: Colors.orange,
                      ),
                    ),
                    labelText: "房間號碼",
                    hintText: "number",
                    labelStyle: TextStyle(color: Colors.orange),
                    icon: Icon(
                      Icons.home,
                      color: Colors.orange,
                    ),
                  ),
                  keyboardType: TextInputType.number,
                ),
              ),
              Divider(height: 30.0),
              Row(children: <Widget>[
                Flexible(
                  child: RadioListTile<String>(
                    value: "自定義",
                    title: Text("自定義"),
                    groupValue: _selectNum,
                    onChanged: (value) {
                      setState(() {
                        _selectNum = value;
                        print(mapPlayerSelect);
                      });
                    },
                  ),
                ),
              ]),
              Row(
                children: <Widget>[
                  Flexible(
                    child: RadioListTile<String>(
                      value: "4人",
                      title: Text("4人"),
                      groupValue: _selectNum,
                      onChanged: (value) {
                        setState(() {
                          _selectNum = value;
                          mapPlayerSelect["Human"]=2;
                          mapPlayerSelect["Wolf"]=1;
                          mapPlayerSelect["Prophesier"]=1;
                          _sliderHuman=2;
                          _sliderWolf=1;
                          _sliderProphesier=1;
                        });
                      },
                    ),
                  ),
                  Flexible(
                    child: RadioListTile<String>(
                      value: "5人",
                      title: Text("5人"),
                      groupValue: _selectNum,
                      onChanged: (value) {
                        setState(() {
                          _selectNum = value;
                          mapPlayerSelect["Human"]=2;
                          mapPlayerSelect["Wolf"]=1;
                          mapPlayerSelect["Prophesier"]=1;
                          mapPlayerSelect["Witch"]=1;
                          _sliderHuman=2;
                          _sliderWolf=1;
                          _sliderProphesier=1;
                          _sliderWitch=1;
                        });
                      },
                    ),
                  ),
                  Flexible(
                    child: RadioListTile<String>(
                      value: "6人",
                      title: Text("6人"),
                      groupValue: _selectNum,
                      onChanged: (value) {
                        setState(() {
                          _selectNum = value;
                          mapPlayerSelect["Human"]=2;
                          mapPlayerSelect["Wolf"]=2;
                          mapPlayerSelect["Prophesier"]=1;
                          mapPlayerSelect["Witch"]=1;
                          _sliderHuman=2;
                          _sliderWolf=2;
                          _sliderProphesier=1;
                          _sliderWitch=1;
                        });
                      },
                    ),
                  ),
                ],
              ),
              Divider(height: 30.0),
              SingleChildScrollView(
                child: Column(
                  children: <Widget>[
                    Row(
                      children: <Widget>[
                        Expanded(
                          child: Slider(
                            activeColor: Colors.pink,
                            min: 0.0,
                            max: 10.0,
                            onChanged: (newRating) {
                              setState(() {
                                _sliderHuman = newRating.toInt();
                                mapPlayerSelect["Human"] = _sliderHuman;
                                _selectNum = "自定義";
                              });
                            },
                            value: _sliderHuman.toDouble(),
                          ),
                          flex: 7,
                        ),
                        Expanded(
                          child: Text("平民"),
                          flex: 1,
                        ),
                        Expanded(
                          child: Text("$_sliderHuman"),
                          flex: 1,
                        ),
                      ],
                    ),
                    Row(
                      children: <Widget>[
                        Expanded(
                          child: Slider(
                            activeColor: Colors.pink,
                            min: 0.0,
                            max: 10.0,
                            onChanged: (newRating) {
                              setState(() {
                                _sliderWolf = newRating.toInt();
                                mapPlayerSelect["Wolf"] = _sliderWolf;
                              });
                            },
                            value: _sliderWolf.toDouble(),
                          ),
                          flex: 7,
                        ),
                        Expanded(
                          child: Text("狼人"),
                          flex: 1,
                        ),
                        Expanded(
                          child: Text("$_sliderWolf"),
                          flex: 1,
                        ),
                      ],
                    ),
                    Row(
                      children: <Widget>[
                        Expanded(
                          child: Slider(
                            activeColor: Colors.pink,
                            min: 0.0,
                            max: 10.0,
                            onChanged: (newRating) {
                              setState(() {
                                _sliderKnight = newRating.toInt();
                                mapPlayerSelect["Knight"] = _sliderKnight;
                              });
                            },
                            value: _sliderKnight.toDouble(),
                          ),
                          flex: 7,
                        ),
                        Expanded(
                          child: Text("騎士"),
                          flex: 1,
                        ),
                        Expanded(
                          child: Text("$_sliderKnight"),
                          flex: 1,
                        ),
                      ],
                    ),
                    Row(
                      children: <Widget>[
                        Expanded(
                          child: Slider(
                            activeColor: Colors.pink,
                            min: 0.0,
                            max: 10.0,
                            onChanged: (newRating) {
                              setState(() {
                                _sliderProphesier = newRating.toInt();
                                mapPlayerSelect["Prophesier"] = _sliderProphesier;
                              });
                            },
                            value: _sliderProphesier.toDouble(),
                          ),
                          flex: 7,
                        ),
                        Expanded(
                          child: Text("預言家"),
                          flex: 1,
                        ),
                        Expanded(
                          child: Text("$_sliderProphesier"),
                          flex: 1,
                        ),
                      ],
                    ),
                    Row(
                      children: <Widget>[
                        Expanded(
                          child: Slider(
                            activeColor: Colors.pink,
                            min: 0.0,
                            max: 10.0,
                            onChanged: (newRating) {
                              setState(() {
                                _sliderHunter = newRating.toInt();
                                mapPlayerSelect["Hunter"] = _sliderHunter;
                              });
                            },
                            value: _sliderHunter.toDouble(),
                          ),
                          flex: 7,
                        ),
                        Expanded(
                          child: Text("獵人"),
                          flex: 1,
                        ),
                        Expanded(
                          child: Text("$_sliderHunter"),
                          flex: 1,
                        ),
                      ],
                    ),
                    Row(
                      
                      children: <Widget>[
                        Expanded(
                          child: Slider(
                            activeColor: Colors.pink,
                            min: 0.0,
                            max: 10.0,
                            onChanged: (newRating) {
                              setState(() {
                                _sliderWitch = newRating.toInt();
                                mapPlayerSelect["Witch"] = _sliderWitch;
                              });
                            },
                            value: _sliderWitch.toDouble(),
                          ),
                          flex: 7,
                        ),
                        Expanded(
                          child: Text("女巫"),
                          flex: 1,
                        ),
                        Expanded(
                          child: Text("$_sliderWitch"),
                          flex: 1,
                        ),
                      ],
                    ),
                  ],
                ),
              ),
              Divider(height: 50.0),
              RaisedButton(
                padding: EdgeInsets.all(10.0),
                color: Colors.orange,
                textColor: Colors.white,
                child: Text("建立"),
                shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(5.0)),
                onPressed: () {
                  // int token = int.parse(_token_input.text);
                  // _websocket_channel = new IOWebSocketChannel.connect(
                  //     ws + host + gamename + "?token=" + _token_input.text);
                  // _websocket_channel.sink.add(token);
                  Navigator.pushNamed(context, "game");
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class GamePage extends StatefulWidget {
  IOWebSocketChannel _websocket_channel;

  String host = "127.0.0.1";
  String gamestart = "/wf/";
  String gamename = "/wf/game";
  String ws = "ws://";
  String wss = "wss://";

  @override
  _GamePageState createState() => _GamePageState();
}

class _GamePageState extends State<GamePage> {
  IOWebSocketChannel _websocket_channel;
  String host = "127.0.0.1";
  String gamename = "wf/game";
  String ws = "ws://";
  String wss = "wss://";

  @override
  void initState() {
    super.initState();
    _phone_speak("遊戲開始");
  }
  

  @override
  Widget build(BuildContext context) {
    TextEditingController _num_input = new TextEditingController();
    _websocket_channel = new IOWebSocketChannel.connect(
        'ws://127.0.0.1:8000/wf/game?token=' + '9');
    return Scaffold(
      appBar: AppBar(
        title: Text(
          "遊戲",
          style: TextStyle(color: Colors.white),
        ),
        backgroundColor: Colors.black87,
      ),
      body: Container(
        child: Column(
          children: <Widget>[
            Padding(
              padding: EdgeInsets.symmetric(
                vertical: 16.0,
                horizontal: 32.0,
              ),
              child: TextField(
                controller: _num_input,
                decoration: InputDecoration(
                  focusedBorder: UnderlineInputBorder(
                    borderSide: BorderSide(
                      color: Colors.orange,
                    ),
                  ),
                  labelText: "輸入號碼",
                  hintText: "ex:1",
                  labelStyle: TextStyle(color: Colors.orange),
                  icon: Icon(
                    Icons.home,
                    color: Colors.orange,
                  ),
                ),
              ),
            ),
            RaisedButton(
              padding: EdgeInsets.all(10.0),
              color: Colors.orange,
              textColor: Colors.white,
              child: Text("確認"),
              shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(5.0)),
              onPressed: () {
                _websocket_channel.sink.add(_num_input.text);
              },
            ),
            Padding(
              padding: EdgeInsets.all(20.0),
            ),
            StreamBuilder(
              stream: _websocket_channel.stream,
              builder: (context, snapshot) {
                String jsonData = snapshot.data;
                switch (jsonData) {
                  case "ping":
                    print("case ping");
                    return Padding(
                      padding: EdgeInsets.all(10.0),
                      child: Container(
                        color: Colors.red,
                        height: 100.0,
                        width: 100.0,
                      ),
                    );
                    break;
                  default:
                    print("default");
                }
                print(jsonData);
                // final jsonMap = json.decode(jsonData);
                // SelectUser selectUser = SelectUser.fromJson(jsonMap);
                // TtsHelper.instance
                //     .setLanguageAndSpeak('${selectUser.sound}', "zh");
                return Padding(
                  padding: EdgeInsets.all(10.0),
                );
                // return Padding(
                //   padding: const EdgeInsets.symmetric(vertical: 24.0),
                //   child: Text(snapshot.hasData ? '${snapshot.data}' : ''),
                // );
              },
            )
          ],
        ),
      ),
    );
  }
}

// _phone_speak 手機說話
void _phone_speak(String sound) {
  TtsHelper.instance.setLanguageAndSpeak(sound, "zh");
}
