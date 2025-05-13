#ifndef CLIENT_H
#define CLIENT_H
#include "creds.h"
#include "eventData.h"
#include <QNetworkAccessManager>
#include <QNetworkRequest>
#include <QNetworkReply>
#include <QJsonDocument>
#include <QJsonObject>
#include <QJsonArray>
#include <QEventLoop>

class Client
{
    QNetworkAccessManager am;
    QString host;
    QString port;
public:
    Client(QString host, QString port);
    bool login(credentials, QString*);
    bool registration(credentials_reg);
    bool update(credentials_reg, QString);
    credentials_reg getUserData(const QString &uid);

    bool addEvent(const EventData &event, const QString &uid);
    QVector<EventData> getEventsByDay(QDate day, const QString uid);
    bool deleteEvent(QString uid,QString eventID);

    bool sendMessage(QString eventID, QString message, QString uid);
    QVector<Message> getMessages(QString eventID);
};

#endif // CLIENT_H
