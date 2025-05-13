#ifndef EVENTDATA_H
#define EVENTDATA_H
#include <QString>
#include <QLineEdit>
#include <QListView>
#include <QDateTime>
#include <QStringList>
struct EventData {

    enum Priority { High = 3, Medium = 2, Low = 1 };

    QString ID;
    QDateTime start;
    QDateTime end;
    QString name;
    QString type;
    Priority prior;
    QString desc;
    QString master;
    QStringList parts;

};

struct Message {
    QString sender;
    QString content;
};

#endif // EVENTDATA_H
