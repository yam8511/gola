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
    Navigator.pushNamed(context, "home");
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
              Column(
                children: <Widget>[
                  Slider(
                    activeColor: Colors.pink,
                    min: 0.0,
                    max: 10.0,
                    onChanged: (newRating) {
                      setState(() {
                        _sliderValue = newRating.toInt();
                        print(_sliderValue);
                      });
                    },
                    value: _sliderValue.toDouble(),
                  ),
                ],
              ),Text(_sliderValue.toString(),)
              
            ],
          ),
        ),
      ),
    );
  }
}

class TriangleClipper extends CustomClipper<Path> {
  @override
  Path getClip(Size size) {
    final path = Path();
    path.lineTo(size.width, 0.0);
    path.lineTo(size.width / 2, size.height);
    path.close();
    return path;
  }

  @override
  bool shouldReclip(TriangleClipper oldClipper) => false;
}

class SettingPage extends StatefulWidget {
  @override
  _SettingPageState createState() => _SettingPageState();
}

class _SettingPageState extends State<SettingPage> {
  TextEditingController _token_input = new TextEditingController();

  IOWebSocketChannel _websocket_channel;

  String host = "127.0.0.1";
  String gamestart = "/wf/";
  String gamename = "/wf/game";
  String ws = "ws://";
  String wss = "wss://";

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
                    labelText: "建立房間",
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
              RaisedButton(
                padding: EdgeInsets.all(10.0),
                color: Colors.orange,
                textColor: Colors.white,
                child: Text("建立"),
                shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(5.0)),
                onPressed: () {
                  // int token = int.parse(_token_input.text);
                  _websocket_channel = new IOWebSocketChannel.connect(
                      ws + host + gamename + "?token=" + _token_input.text);

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

////////////////////////////////
class MyHomePage extends StatefulWidget {
  MyHomePage({Key key, this.title}) : super(key: key);

  final String title;

  @override
  _MyHomePageState createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {
  int _counter = 0;

  void _incrementCounter() {
    setState(() {
      _counter++;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            Text(
              'You have pushed the button this many times:',
            ),
            Text(
              '$_counter',
              style: Theme.of(context).textTheme.display1,
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: _incrementCounter,
        tooltip: 'Increment',
        child: Icon(Icons.add),
      ),
    );
  }
}

void _phone_speak(String sound) {
  TtsHelper.instance.setLanguageAndSpeak(sound, "zh");
}
