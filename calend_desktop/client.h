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
#include <QHttpMultiPart>
#include <QHttpPart>
#include <QFile>
#include <QFileInfo>
#include <QPair>
#include <QList>
#include <QObject>
#include <QMetaObject>
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

    bool addEvent(const EventData &event, const QString &uid, QString* eventID = nullptr);
    QVector<EventData> getEventsByDay(QDate day, const QString uid);
    bool deleteEvent(QString uid,QString eventID);

    bool sendMessage(QString eventID, QString message, QString uid);
    QVector<Message> getMessages(QString eventID);

    bool uploadAttachment(const QString& eventID, QFile& file);
    QList<QPair<QString, QString>> getAttachments(const QString& eventID);

    QVector<EventData> getAllUserEvents(const QString &uid);

    bool updateParticipation(const QString& eventID, const QString& uid, bool state);

    bool scheduleNotification(const QString& eventID, const QDateTime& deadline);

    EventData getEventByID(const QString& eventID);

    bool updateEvent(const EventData &event);

    QVector<Participant> getParticipants(const QString& eventID);
};

#endif // CLIENT_H
