#ifndef EVENT_H
#define EVENT_H
#include <QString>
#include <QWidget>
#include <QDateTime>
#include <QStringList>
#include "cfg.h"
#include "client.h"
#include "eventData.h"
#include <QLineEdit>
#include <QListWidget>
#include <QMessageBox>

namespace Ui {
class Event;
}

class Event : public QWidget
{
    Q_OBJECT

public:
    explicit Event(QString uid, QWidget *parent = nullptr);
    ~Event();
    void fillFormWithEventData(const EventData& data);
signals:
    void eventUpdated();

private slots:

    void on_create_event_clicked();

    QDateTime calculateNotificationTime(const QDateTime& eventStart, int notifyIndex);

private:
    Ui::Event *ui;
    Client cli;
    QString uid;
    bool isEditMode = false;
    QString originalEventId;
};

#endif // EVENT_H
