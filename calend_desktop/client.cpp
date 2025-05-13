#include "client.h"

Client::Client(QString host, QString port) : am(), host(host), port(port) {
}

bool Client::login(credentials creds, QString* uid) {
    QNetworkRequest req(QUrl("http://"+host+":"+port+"/users/login"));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    QJsonObject json;
    json["mail"] = creds.mail;
    json["pass"] = creds.pass;
    QByteArray data = QJsonDocument(json).toJson();
    QNetworkReply *reply = am.post(req, data);
    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();
    if (reply->error() != QNetworkReply::NoError) {
        qDebug() << "Error:" << reply->errorString();
        return false;
    } else {
        QJsonDocument jsonresponse = QJsonDocument::fromJson(reply->readAll());
        *uid = jsonresponse["uid"].toString();
        return true;
    }
}
bool Client::registration(credentials_reg creds_reg){
    QNetworkRequest req(QUrl("http://"+host+":"+port+"/users/register"));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    QJsonObject json;
    json["mail"] = creds_reg.mail;
    json["pass"] = creds_reg.pass;
    json["f_name"] = creds_reg.f_name;
    json["s_name"] = creds_reg.s_name;
    json["t_name"] = creds_reg.t_name;
    json["dep"] = creds_reg.department;
    json["pos"] = creds_reg.pos;

    QByteArray data = QJsonDocument(json).toJson();
    QNetworkReply *reply = am.post(req, data);

    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();

    if (reply->error() != QNetworkReply::NoError) {
        qDebug() << "Error:" << reply->errorString();
        return false;
    }
    return true;

}

bool Client::update(credentials_reg creds, QString uid) {
    QNetworkRequest req(QUrl("http://"+host+":"+port+"/users/"+uid+"/update"));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    QJsonObject json;
    json["mail"] = creds.mail;
    json["pass"] = creds.pass;
    json["f_name"] = creds.f_name;
    json["s_name"] = creds.s_name;
    json["t_name"] = creds.t_name;
    json["dep"] = creds.department;
    json["pos"] = creds.pos;

    QByteArray data = QJsonDocument(json).toJson();
    QNetworkReply *reply = am.post(req, data);

    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();

    if (reply->error() != QNetworkReply::NoError) {
        qDebug() << "Error:" << reply->errorString();
        return false;
    }
    return true;

}

credentials_reg Client::getUserData(const QString &uid) {
    QNetworkRequest req(QUrl("http://"+host+":"+port+"/users/"+uid+"/profile"));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    QNetworkReply *reply = am.get(req);

    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();
    credentials_reg out;
    if (reply->error() != QNetworkReply::NoError) {
        qDebug() << "Error:" << reply->errorString();
        return out;
    }
    QByteArray raw = reply->readAll();
    QJsonDocument jsonresponse = QJsonDocument::fromJson(raw);
    out.mail = jsonresponse["mail"].toString();
    out.f_name = jsonresponse["f_name"].toString();
    out.s_name = jsonresponse["s_name"].toString();
    out.t_name = jsonresponse["t_name"].toString();
    out.department = jsonresponse["dep"].toString();
    out.pos = jsonresponse["pos"].toString();
    reply->deleteLater();
    return out;
}

bool Client::addEvent(const EventData &event, const QString &uid)
{
    QNetworkRequest req(QUrl("http://" + host + ":" + port + "/events/add"));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");

    QJsonObject json;
    json["master"] = uid;
    json["name"] = event.name;
    json["desc"] = event.desc;
    json["type"] = event.type;
    json["prior"] = event.prior;  // 0 1 или 2
    json["start"] = event.start.toUTC().toString(Qt::ISODateWithMs);
    json["end"] = event.end.toUTC().toString(Qt::ISODateWithMs);

    QJsonArray partsArray;
    for (const QString &email : event.parts) {
        partsArray.append(email);
    }
    json["parts"] = partsArray;

    QByteArray data = QJsonDocument(json).toJson();
    QNetworkReply *reply = am.post(req, data);

    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();

    if (reply->error() != QNetworkReply::NoError) {
        qDebug() << "Ошибка при создании события:" << reply->errorString();
        return false;
    }

    return true;
}

QVector<EventData> Client::getEventsByDay(QDate day, const QString uid) {
    QVector<EventData> result;
    QNetworkRequest req(QUrl("http://"+host+":"+port+"/events/"+uid+"/day?day="+day.toString(Qt::ISODate)));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    QNetworkReply *reply = am.get(req);

    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();
    if (reply->error() != QNetworkReply::NoError) {
        qDebug() << "Error:" << reply->errorString();
        return result;
    }
    QByteArray raw = reply->readAll();
    QJsonParseError parseError;
    QJsonDocument doc = QJsonDocument::fromJson(raw, &parseError);
    if (doc.isNull()) {
        qDebug() << "error: " << parseError.errorString() << "at: " << parseError.offset;

        int start = qMax(0, parseError.offset);
        int length = 40;
        qDebug() << "Error context:"
                 << raw.mid(start, length);
        return result;
    }
    QJsonArray events = doc["events"].toArray();
    for(int i = 0; i < events.size(); i++) {
        const QJsonObject event = events.at(i).toObject();
        EventData item;
        item.prior = static_cast<EventData::Priority>(event["prior"].toInt());
        for (const QJsonValue& ev : event["parts"].toArray()) {
            item.parts.append(ev.toString());
        }
        item.ID = event["id"].toString();
        item.master = event["master"].toString();
        item.desc = event["desc"].toString();
        item.name = event["name"].toString();
        item.type = event["type"].toString();
        item.start = QDateTime::fromString(event["start"].toString(), Qt::ISODateWithMs);
        item.end = QDateTime::fromString(event["end"].toString(), Qt::ISODateWithMs);
        result.append(item);
    }
    return result;
}

bool Client::deleteEvent(QString uid, QString eventID) {
    QNetworkRequest req(QUrl("http://"+host+":"+port+"/events/"+uid+"/delete?id="+eventID));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    QNetworkReply *reply = am.deleteResource(req);
    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();
    if (reply->error() != QNetworkReply::NoError) {
        qDebug() << "Error:" << reply->errorString();
        return false;
    }
    return true;
}

bool Client::sendMessage(QString eventID, QString message, QString uid) {
    QNetworkRequest req(QUrl("http://"+host+":"+port+"/chats/"+eventID));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");

    QJsonObject body;
    body["cont"] = message;
    body["sender"] = uid;
    QByteArray rawbody = QJsonDocument(body).toJson();
    QNetworkReply *reply = am.post(req, rawbody);
    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();
    if (reply->error() != QNetworkReply::NoError) {
        qDebug() << "Error:" << reply->errorString();
        return false;
    }
    return true;
}

QVector<Message> Client::getMessages(QString eventID) {
    QNetworkRequest req(QUrl("http://"+host+":"+port+"/chats/"+eventID));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    QVector<Message> result;
    QNetworkReply *reply = am.get(req);

    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();
    if (reply->error() != QNetworkReply::NoError) {
        qDebug() << "Error:" << reply->errorString();
        return result;
    }

    QByteArray raw = reply->readAll();
    QJsonParseError parseError;
    QJsonDocument doc = QJsonDocument::fromJson(raw, &parseError);
    if (doc.isNull()) {
        qDebug() << "error: " << parseError.errorString() << "at: " << parseError.offset;

        int start = qMax(0, parseError.offset);
        int length = 40;
        qDebug() << "Error context:"
                 << raw.mid(start, length);
        return result;
    }
    for (const QJsonValue& v : doc["messages"].toArray()) {
        Message msg;
        msg.content = v.toObject()["cont"].toString();
        msg.sender = v.toObject()["sender"].toString();
        result.append(msg);
    }
    return result;
}
