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


namespace Ui {
class Event;
}

class Event : public QWidget
{
    Q_OBJECT

public:
    explicit Event(QString uid, QWidget *parent = nullptr);
    ~Event();

private slots:

    void on_create_event_clicked();

private:
    Ui::Event *ui;
    Client cli;
    QString uid;
};

#endif // EVENT_H
