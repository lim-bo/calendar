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

bool Client::addEvent(const EventData &event, const QString &uid, QString* eventID)
{
    QNetworkRequest req(QUrl("http://" + host + ":" + port + "/events/add"));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");

    QJsonObject json;
    json["master"] = uid;
    json["name"] = event.name;
    json["desc"] = event.desc;
    json["type"] = event.type;
    json["prior"] = event.prior;  // 1 2 или 3
    json["start"] = event.start.toUTC().toString(Qt::ISODateWithMs);
    json["end"] = event.end.toUTC().toString(Qt::ISODateWithMs);

    QJsonArray partsArray;
    for (const Participant &participant : event.parts) {
        partsArray.append(participant.uid);
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
    *eventID = QString::fromUtf8(reply->rawHeader("eventID"));
    return true;
}

bool Client::scheduleNotification(const QString& eventID, const QDateTime& deadline)
{

    QNetworkRequest req(QUrl("http://" + host + ":" + port + "/events/" + eventID + "/notify?deadline=" + deadline.toString("yyyy-MM-dd_HH:mm:ss")));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");

    QNetworkReply *reply = am.post(req, QByteArray());
    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();

    return reply->error() == QNetworkReply::NoError;
}

EventData Client::getEventByID(const QString& eventID) {
    QNetworkRequest req(QUrl("http://" + host + ":" + port + "/events/" + eventID));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    QNetworkReply *reply = am.get(req);
    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();

    if (reply->error() == QNetworkReply::NoError) {
        qDebug() << reply->error();
        return EventData{};
    }

    QByteArray raw = reply->readAll();
    QJsonParseError parseError;
    QJsonDocument event = QJsonDocument::fromJson(raw, &parseError);
    if (event.isNull()) {
        qDebug() << "error: " << parseError.errorString() << "at: " << parseError.offset;

        int start = qMax(0, parseError.offset);
        int length = 40;
        qDebug() << "Error context:"
                 << raw.mid(start, length);
        return EventData{};
    }
    EventData item;
    item.ID = event["id"].toString();
    item.master = event["master"].toString();
    item.desc = event["desc"].toString();
    item.name = event["name"].toString();
    item.type = event["type"].toString();
    item.start = QDateTime::fromString(event["start"].toString(), Qt::ISODateWithMs);
    item.end = QDateTime::fromString(event["end"].toString(), Qt::ISODateWithMs);
    QJsonArray partsArray = event["parts"].toArray();
    for (const QJsonValue& partValue : partsArray) {
        QJsonObject partObj = partValue.toObject();
        Participant participant;
        participant.uid = partObj["uid"].toString();
        participant.accepted = partObj["accepted"].toBool();
        item.parts.append(participant);
    }
    return item;
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

        item.ID = event["id"].toString();
        item.master = event["master"].toString();
        item.desc = event["desc"].toString();
        item.name = event["name"].toString();
        item.type = event["type"].toString();
        item.start = QDateTime::fromString(event["start"].toString(), Qt::ISODateWithMs);
        item.end = QDateTime::fromString(event["end"].toString(), Qt::ISODateWithMs);

        QJsonArray partsArray = event["parts"].toArray();
        for (const QJsonValue& partValue : partsArray) {
            QJsonObject partObj = partValue.toObject();
            Participant participant;
            participant.uid = partObj["uid"].toString();
            participant.accepted = partObj["accepted"].toBool();
            item.parts.append(participant);
        }

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




bool Client::uploadAttachment(const QString& eventID, QFile& file)
{
    file.open(QIODevice::ReadOnly);
    if (!file.isOpen()) {
        qDebug() << "Не получилось открыть файл";
        return false;
    }

    QFileInfo fileInfo(file);
    QString fileName = fileInfo.fileName();

    QNetworkRequest req(QUrl("http://" + host + ":" + port + "/attachs/" + eventID));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    QJsonObject obj;
    QJsonArray data;
    QByteArray fileData = file.readAll();
    for (const auto& item : fileData) {
        data.append(static_cast<unsigned char>(item));
    }
    obj["name"] = fileName;
    obj["data"] = data;
    QNetworkReply *reply = am.post(req, QJsonDocument(obj).toJson());

    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();

    reply->deleteLater();
    file.close();//анекдот: как называются маленькие нервные люди ........................... микроволновки

    return reply->error() == QNetworkReply::NoError;
}

QList<QPair<QString, QString>> Client::getAttachments(const QString& eventID)
{
    QList<QPair<QString, QString>> attachmentsList;
    QNetworkRequest req(QUrl("http://" + host + ":" + port + "/attachs/" + eventID));
    QNetworkReply *reply = am.get(req);

    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();

    if (reply->error() != QNetworkReply::NoError) {
        qDebug() << "Ошибка получения вложений:" << reply->errorString();
        return attachmentsList;
    }

    QJsonArray attachmentsArray = QJsonDocument::fromJson(reply->readAll()).array();

    for (const QJsonValue &value : attachmentsArray) {
        QJsonObject attachment = value.toObject();
        attachmentsList.append(qMakePair(
            attachment["name"].toString(),
            attachment["data"].toString()
            ));
    }

    return attachmentsList;
}

QVector<EventData> Client::getAllUserEvents(const QString &uid) {
    QVector<EventData> result;
    QNetworkRequest req(QUrl("http://"+host+":"+port+"/events/"+uid));//как называются маленькие деньги, которым страшно?..................сущие копейки
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
        return result;
    }

    QJsonArray events = doc["events"].toArray();
    for(int i = 0; i < events.size(); i++) {
        const QJsonObject event = events.at(i).toObject();
        EventData item;
        item.prior = static_cast<EventData::Priority>(event["prior"].toInt());
        item.ID = event["id"].toString();
        item.master = event["master"].toString();
        item.desc = event["desc"].toString();
        item.name = event["name"].toString();
        item.type = event["type"].toString();
        item.start = QDateTime::fromString(event["start"].toString(), Qt::ISODateWithMs);
        item.end = QDateTime::fromString(event["end"].toString(), Qt::ISODateWithMs);

        QJsonArray partsArray = event["parts"].toArray();
        for (const QJsonValue& partValue : partsArray) {
            QJsonObject partObj = partValue.toObject();
            Participant participant;
            participant.uid = partObj["uid"].toString();
            participant.accepted = partObj["accepted"].toBool();
            item.parts.append(participant);
        }

        result.append(item);
    }
    return result;
}


bool Client::updateParticipation(const QString& eventID, const QString& uid, bool state)
{
    QNetworkRequest req(QUrl("http://" + host + ":" + port + "/events/" + eventID + "/" + uid + "?state=" + (state ? "1" : "0")));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");

    QNetworkReply *reply = am.post(req, QByteArray());
    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();

    if (reply->error() != QNetworkReply::NoError) {

        return false;
    }
    return true;
}


bool Client::updateEvent(const EventData &event) {
    QNetworkRequest req(QUrl("http://"+host+":"+port+"/events/update"));
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");


    QJsonObject json;
    json["id"] = event.ID;
    json["master"] = event.master;
    json["name"] = event.name;
    json["desc"] = event.desc;
    json["type"] = event.type;
    json["prior"] = event.prior;
    json["start"] = event.start.toUTC().toString(Qt::ISODateWithMs);
    json["end"] = event.end.toUTC().toString(Qt::ISODateWithMs);
    QJsonArray partsArray;
    for (const Participant &participant : event.parts) {
        partsArray.append(participant.uid);
    }
    json["parts"] = partsArray;
    QByteArray data = QJsonDocument(json).toJson();
    QNetworkReply *reply = am.post(req, data);

    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();
    if (reply->error() != QNetworkReply::NoError) {
        qDebug() << "Ошибка при обновлении события:" << reply->errorString();
        return false;
    }
    return true;
}

QVector<Participant> Client::getParticipants(const QString& eventID) {
    QNetworkRequest req(QUrl("http://"+host+":"+port+"/events/"+eventID+"/parts"));

    QNetworkReply *reply = am.get(req);

    QEventLoop loop;
    QObject::connect(reply, &QNetworkReply::finished, &loop, &QEventLoop::quit);
    loop.exec();

    if (reply->error() != QNetworkReply::NoError) {
        qDebug() << "Ошибка при получении списка участников:" << reply->errorString();
        return QVector<Participant>{};
    }


    QVector<Participant> result;
    QByteArray raw = reply->readAll();
    QJsonParseError parseError;
    QJsonDocument doc = QJsonDocument::fromJson(raw, &parseError);

    if (doc.isNull()) {
        qDebug() << "error: " << parseError.errorString() << "at: " << parseError.offset;
        return result;
    }

    for (const QJsonValue& part : doc["parts"].toArray()) {
        QJsonObject obj = part.toObject();
        Participant item;
        item.uid = obj["email"].toString();
        item.accepted = obj["accepted"].toBool();
        result.append(item);
    }
    return result;
}
